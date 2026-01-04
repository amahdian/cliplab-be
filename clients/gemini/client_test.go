package gemini

import (
	"testing"
	"time"

	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/pkg/logger"
)

func TestAnalyzeVideo_Success(t *testing.T) {

	date, _ := time.Parse("2006-01-02 15:04:05", "2025-11-26 17:01:53")
	client := NewClient("https://generativelanguage.googleapis.com", "AIzaSyC64CVnhAHSADMloYwCe8hbMalSCSnXgJg")
	resp, _ := client.AnalyzeVideo(
		model.PlatformInstagram,
		"https://scontent-prg1-1.cdninstagram.com/o1/v/t2/f2/m86/AQM5-tU7kYwNnMkDKY2gOnuGErhlOQIk4Vya3Px7nAs-u5TGcgYCjt5bEIaOSELdqzjcXFWQ3bRtwUC3MrN_3vKWT-9fFJRLgyeNn-M.mp4?_nc_cat=102&_nc_sid=5e9851&_nc_ht=scontent-prg1-1.cdninstagram.com&_nc_ohc=SSqTmPLde9kQ7kNvwFeqtX1&efg=eyJ2ZW5jb2RlX3RhZyI6Inhwdl9wcm9ncmVzc2l2ZS5JTlNUQUdSQU0uQ0xJUFMuQzMuNzIwLmRhc2hfYmFzZWxpbmVfMV92MSIsInhwdl9hc3NldF9pZCI6NDQ3NTM0NjU4Mjc1MTQ5MSwiYXNzZXRfYWdlX2RheXMiOjM4LCJ2aV91c2VjYXNlX2lkIjoxMDA5OSwiZHVyYXRpb25fcyI6MTc2LCJ1cmxnZW5fc291cmNlIjoid3d3In0%3D&ccb=17-1&_nc_gid=xeo_JY0R_QeB4hoa_8y8vg&_nc_zt=28&vs=15cb1e2f76bc7f5f&_nc_vs=HBksFQIYUmlnX3hwdl9yZWVsc19wZXJtYW5lbnRfc3JfcHJvZC81OTRDMUFDMzY5NDE3MEE0RUY4QUU0REI4N0Y4MDU5Rl92aWRlb19kYXNoaW5pdC5tcDQVAALIARIAFQIYOnBhc3N0aHJvdWdoX2V2ZXJzdG9yZS9HRnN0LVNKZ0pzU3RaVzRHQUFIUWVoejRZTFFDYnN0VEFRQUYVAgLIARIAKAAYABsCiAd1c2Vfb2lsATEScHJvZ3Jlc3NpdmVfcmVjaXBlATEVAAAmhqT-tLqT8w8VAigCQzMsF0BmGuFHrhR7GBJkYXNoX2Jhc2VsaW5lXzFfdjERAHX-B2XmnQEA&oh=00_AfrtMsK0s_djxrWWZ3WxFI6Oz9sKtxV0yGIclz-52Aw7Ow&oe=695B35AF",
		"@sadiesink_ (Sadie Sink) and @maya_hawke (Maya Hawke) are returning to Hawkins for one last adventure 🚲📺⁣\\n ⁣\\nThe duo sat down to talk about the final season of @strangerthingstv (“Stranger Things”), their fav hype songs (s/o to the “Hamilton” soundtrack) and the power of friendship in the latest Close Friends Only: Speed Round 🌟",
		[]string{},
		[]string{
			"now this is a dynamic duo",
			"@vikusja181 ❤️😍💨",
			"@instagram 🔥👏",
			"Max and Robin really need a spinoff of their own!! 😍",
			"@instagram 💖👑",
			"@instagram yg huu",
			"@beforv i heard famous people were talking in comments? ✌️",
			"Two of the bestest bests ❤️❤️❤️",
			"@nellfisher_ .hlw misis",
			"@nellfisher_ like seriously 😍😍",
			"@nellfisher_ ❤️🙌",
			"This duo! ❤️🙌",
			"@instagram",
			"@sadieslaugh yes",
			"@sadieslaugh the Sadie’s stay in close orbit ,_’,’ 😂 @sadieslaugh @sadiesink_ yeah we love the way you n @maya_hawke jive",
			"The fact that Sadie Is incredible at Mario Kart like Max is incredible at videogames 🥹❤️‍🩹",
			"@instagram 🏆💎",
			"@_filevenvideo 😍😍😍😍😍😍",
			"@_filevenvideo اع مطمعننی",
			"🥰🥰🥰🥰🥰🥰🥰🥰🥰🥰🥰",
			"👏👏👏👏👏",
			"👏👏👏👏👏👏",
			"😻😻",
			"❤️❤️❤️",
			"So beautiful 😍 🤩 👌 ❣️ 💖 💗",
			"So nice 💚💚💚❤️❤️❤️❤️🔥🔥🔥😍😍😍",
			"Oh my gosh yes; Medieval Times is amazing!",
			"@instagram 💛😀",
			"❤️❤️❤️",
			"So excited for their return! The friendship between them is everything. Can't wait to see how it all wraps up! 🚲✨",
			"❤️❤️❤️",
		},
		map[string]float64{
			"like_count":      990650,
			"comment_count":   7089,
			"view_count":      26661931,
			"play_count":      113670205,
			"engagement_rate": 0.14,
		},
		map[string]float64{
			"follower_count":          698663862,
			"average_like_count":      631735,
			"average_comment_count":   10771,
			"average_view_count":      0,
			"average_play_count":      85129456,
			"average_engagement_rate": 0.091,
		},
		date,
		"US",
	)

	logger.Debug(resp)

}
