package msg

import "encoding/json"

type MessageLevel int

const (
	Fatal MessageLevel = iota
	Error
	Warning
	Info
)

func (ml MessageLevel) String() string {
	switch ml {
	case Fatal:
		return "Fatal"
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	default:
		return "Info"
	}
}

func ParseMessageLevel(str string) MessageLevel {
	switch str {
	case "Fatal":
		return Fatal
	case "Error":
		return Error
	case "Warning":
		return Warning
	default:
		return Info
	}
}

type Message struct {
	Text  string
	Level MessageLevel
}

type MessageGroup struct {
	Group    string     `json:"group"`
	Messages []*Message `json:"messages"`
	Order    int        `json:"-"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	message := map[string]string{
		"text":  MakePlainText(m.Text),
		"level": m.Level.String(),
	}
	return json.Marshal(message)
}
func (m *Message) UnmarshalJSON(data []byte) error {
	var j map[string]string
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	m.Text = j["text"]
	m.Level = ParseMessageLevel(j["level"])
	return nil
}
