package helix

import (
	"net/http"
	"testing"
)

func TestGetUsers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		IDs         []string
		Logins      []string
		respBody    string
		expectUsers []string
	}{
		{
			http.StatusBadRequest,
			[]string{},
			[]string{},
			`{"error":"Bad Request","status":400,"message":"Must provide an ID, Login or OAuth Token"}`,
			[]string{},
		},
		{
			http.StatusOK,
			[]string{"26301881"},
			[]string{"summit1g"},
			`{"data":[{"id":"26301881","login":"sodapoppin","display_name":"sodapoppin","type":"","broadcaster_type":"partner","description":"Wtf do i write here? Click my stream, or i scream.","profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/sodapoppin-profile_image-10049b6200f90c14-300x300.png","offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/sodapoppin-channel_offline_image-2040c6fcacec48db-1920x1080.jpeg","view_count":190154823},{"id":"26490481","login":"summit1g","display_name":"summit1g","type":"","broadcaster_type":"partner","description":"I'm a competitive CounterStrike player who likes to play casually now and many other games. You will mostly see me play CS, H1Z1,and single player games at night. There will be many othergames played on this stream in the future as they come out:D","profile_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/200cea12142f2384-profile_image-300x300.png","offline_image_url":"https://static-cdn.jtvnw.net/jtv_user_pictures/summit1g-channel_offline_image-e2f9a1df9e695ec1-1920x1080.png","view_count":202707885}]}`,
			[]string{"sodapoppin", "summit1g"},
		},
	}

	for _, testCase := range testCases {
		c := newMockClient("cid", newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetUsers(&UsersParams{
			IDs:    testCase.IDs,
			Logins: testCase.Logins,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be %s, got %s", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusBadRequest {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "Must provide an ID, Login or OAuth Token"
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if resp.Data.Users[0].Login != testCase.expectUsers[0] { // sodapoppin
			t.Errorf("expected username 1 to be %s, got %s", testCase.expectUsers[0], resp.Data.Users[0].Login)
		}

		if resp.Data.Users[1].Login != testCase.expectUsers[1] { // summit1g
			t.Errorf("expected username 2 to be %s, got %s", testCase.expectUsers[0], resp.Data.Users[0].Login)
		}
	}
}
