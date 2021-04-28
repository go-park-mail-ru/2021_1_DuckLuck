package sanitizer

import "github.com/microcosm-cc/bluemonday"

type Sanitizer struct {
	policy *bluemonday.Policy
}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		policy: bluemonday.UGCPolicy(),
	}
}

func (s *Sanitizer) Sanitize(text string) string {
	return s.policy.Sanitize(text)
}
