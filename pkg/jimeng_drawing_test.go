package pkg

import "testing"

func TestJimengDrawing(t *testing.T) {
	imageURL := "https://uranus-houhou.oss-cn-beijing.aliyuncs.com/uranus/user/images/2026/03/23/20260323130839_501e5bfc-e014-4886-b10f-a92b0a1e31c3.jpg"
	config := JimengConfig{
		Enabled:   true,
		SessionID: []string{"xxxxxxxx"},
		JimengRequest: JimengRequest{
			Prompt: "将车移到停车位内",
			Images: []*string{&imageURL},
		},
	}

	imageURLs, err := JimengImageCompositions(&config)
	if err != nil {
		t.Fatalf("JimengImageCompositions failed: %v", err)
	}
	t.Logf("Generated image URLs: %v", imageURLs)
}
