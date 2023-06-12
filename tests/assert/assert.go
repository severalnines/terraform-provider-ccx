package assert

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	testhttp "github.com/severalnines/terraform-provider-ccx/tests/http"
)

var replacer = strings.NewReplacer(
	"\r", "[R]", "\n", "[N]",
)

// Strings assertion
func Strings(t *testing.T, name string, got, want string) bool {
	if got == want {
		return true
	}

	t.Errorf("%s not the same.\nGot  = %s\nWant = %s\n", name, replacer.Replace(got), replacer.Replace(want))

	return false
}

func Errors(t *testing.T, got, want error) bool {
	if (got == nil) && (want != nil) {
		t.Errorf("Got error is nil, want error is not nil = %s\n", want.Error())
		return false
	} else if (want == nil) && (got != nil) {
		t.Errorf("Want error is nil, got error is not nil = %s\n", got.Error())
		return false
	}

	if errors.Is(got, want) {
		return true
	}

	if got.Error() == want.Error() {
		return true
	}

	t.Errorf("Got error =  %s, want error = %s\n", got.Error(), want.Error())
	return false
}

func HttpHeaders(t *testing.T, name string, got, want http.Header) bool {
	same := true

	if len(got) != len(want) {
		t.Errorf("%s not the same.\n Length Got  = %d\nLength Want = %d\n", name, len(got), len(want))

		same = false
	}

	for k := range got {
		if _, ok := want[k]; !ok {
			t.Errorf("got unwanted key: %s\n", k)
			same = false
		}
	}

	for k := range want {
		if _, ok := got[k]; !ok {
			t.Errorf("did not get wanted key: %s\n", k)
			same = false
		}
	}

	return same
}

func TestRequests(t *testing.T, name string, got, want testhttp.Request) bool {
	same := true

	if !Strings(t, name+".Host", got.Host, want.Host) {
		same = false
	}

	if !Strings(t, name+".Path", got.Path, want.Path) {
		same = false
	}

	if !Strings(t, name+".Query", got.Query, want.Query) {
		same = false
	}

	if !Strings(t, name+".Method", got.Method, want.Method) {
		same = false
	}

	if !Strings(t, name+".Body", got.Body, want.Body) {
		same = false
	}

	if !HttpHeaders(t, name+".Headers", got.Header, want.Header) {
		same = false
	}

	return same
}

func TestRequestsAll(t *testing.T, name string, got, want []testhttp.Request) bool {
	same := true
	if len(got) != len(want) {
		t.Errorf("Got len(%s)  = %d\nWant len(%s) = %d\n", name, len(got), name, len(want))
		return false
	}

	for i := range got {
		same = same && TestRequests(t, fmt.Sprintf("%s[%d]", name, i), got[i], want[i])
	}

	return same
}
