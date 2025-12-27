package pg

import (
	"fmt"
	"sync"

	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/storage"
	"gorm.io/gorm"
)

type ormTxn struct {
	id int

	txn      *gorm.DB
	resolved bool

	parent *ormTxn
	next   *ormTxn
}

func newOrmTxn(db *gorm.DB) (*ormTxn, error) {
	txn := db.Begin()
	if db.Error != nil {
		return nil, db.Error
	}
	return &ormTxn{
		id:  0,
		txn: txn,
	}, nil
}

func newNestOrmTxn(cur *ormTxn) (*ormTxn, error) {
	id := cur.id + 1
	sp := fmt.Sprintf("sp%d", id)
	txn := cur.txn.SavePoint(sp)
	if cur.txn.Error != nil {
		return nil, cur.txn.Error
	}
	next := &ormTxn{
		id:     id,
		txn:    txn,
		parent: cur,
	}
	cur.next = next
	return next, nil
}

type ormSession struct {
	db *gorm.DB

	cur *ormTxn

	mu sync.Mutex
}

func (ses *ormSession) Begin() (storage.Session, error) {
	ses.mu.Lock()
	defer ses.mu.Unlock()

	if ses.cur == nil {
		cur, err := newOrmTxn(ses.db)
		if err != nil {
			return nil, err
		}
		return &ormSession{
			db:  cur.txn,
			cur: cur,
		}, nil
	}
	cur, err := newNestOrmTxn(ses.cur)
	if err != nil {
		return nil, err
	}
	return &ormSession{
		db:  cur.txn,
		cur: cur,
	}, nil
}

func (ses *ormSession) Rollback() error {
	ses.mu.Lock()
	defer ses.mu.Unlock()
	cur := ses.cur
	if cur.next != nil && !cur.next.resolved {
		return errs.Newf(errs.Internal, nil, "Next txn needs to be resolved first.")
	}
	if cur.resolved {
		return nil
	}
	cur.resolved = true

	if ses.cur.parent == nil {
		ses.db.Rollback()
		if ses.db.Error != nil {
			return ses.db.Error
		}
		return nil
	}
	sp := fmt.Sprintf("sp%d", ses.cur.id)
	ses.db.RollbackTo(sp)
	if ses.db.Error != nil {
		return ses.db.Error
	}
	return nil
}

func (ses *ormSession) Commit() error {
	ses.mu.Lock()
	defer ses.mu.Unlock()
	cur := ses.cur
	if cur.next != nil && !cur.next.resolved {
		return errs.Newf(errs.Internal, nil, "Next txn needs to be resolved first.")
	}
	if cur.resolved {
		return nil
	}
	cur.resolved = true

	if ses.cur.parent != nil {
		// do nothing, gorm does not provide a way to release the savepoint.
		// we just pop out the last element of the savepoint by decrement the
		// right boundary of savepoints
		return nil
	}
	ses.db.Commit()
	if ses.db.Error != nil {
		return ses.db.Error
	}
	return nil
}

func (ses *ormSession) Close() error {
	return nil
}

var _ storage.Session = &ormSession{}
