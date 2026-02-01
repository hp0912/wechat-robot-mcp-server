package pkg

import "testing"

func TestZImageDrawing(t *testing.T) {
	imageURLs, err := ZImageDrawing(&ZImageConfig{
		BaseURL: "https://api-inference.modelscope.cn/",
		ApiKey:  "xxxxx",
		Prompt:  "一只戴着眼镜的猫，卡通风格",
	})
	if err != nil {
		t.Fatalf("ZImageDrawing failed: %v", err)
	}
	if len(imageURLs) == 0 {
		t.Fatal("No images returned")
	}
	for _, url := range imageURLs {
		if url != nil {
			t.Logf("Generated image URL: %s", *url)
		} else {
			t.Log("Generated image URL is nil")
		}
	}
}
