package utilhtml

import "testing"

func TestSanitize(t *testing.T) {
	type args struct {
		unsafeHTML string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test 1", args{`<script>alert('XSS')</script><p onclick="steal()">Hello, <b>World</b>!</p>`},
			`<p>Hello, <b>World</b>!</p>`, false},
		{"test 2", args{`<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>`},
			`XSS`, false},
		{"test 3", args{`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`},
			`<a href="http://www.google.com" rel="nofollow">Google</a>`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sanitize(tt.args.unsafeHTML)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sanitize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}
