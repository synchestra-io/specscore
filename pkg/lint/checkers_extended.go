package lint

// Extended checkers are stubs in MVP - will be implemented in Phase 2.

// headingLevelsChecker verifies consistent heading level structure.
type headingLevelsChecker struct{}

func newHeadingLevelsChecker() checker {
	return &headingLevelsChecker{}
}

func (c *headingLevelsChecker) name() string     { return "heading-levels" }
func (c *headingLevelsChecker) severity() string { return "warning" }

func (c *headingLevelsChecker) check(_ string) ([]Violation, error) {
	// TODO: Phase 2
	return []Violation{}, nil
}

// featureRefSyntaxChecker validates feature reference syntax.
type featureRefSyntaxChecker struct{}

func newFeatureRefSyntaxChecker() checker {
	return &featureRefSyntaxChecker{}
}

func (c *featureRefSyntaxChecker) name() string     { return "feature-ref-syntax" }
func (c *featureRefSyntaxChecker) severity() string { return "error" }

func (c *featureRefSyntaxChecker) check(_ string) ([]Violation, error) {
	// TODO: Phase 2
	return []Violation{}, nil
}

// internalLinksChecker validates internal link resolution.
type internalLinksChecker struct{}

func newInternalLinksChecker() checker {
	return &internalLinksChecker{}
}

func (c *internalLinksChecker) name() string     { return "internal-links" }
func (c *internalLinksChecker) severity() string { return "error" }

func (c *internalLinksChecker) check(_ string) ([]Violation, error) {
	// TODO: Phase 2
	return []Violation{}, nil
}

// forwardRefsChecker validates no references to non-existent features.
type forwardRefsChecker struct{}

func newForwardRefsChecker() checker {
	return &forwardRefsChecker{}
}

func (c *forwardRefsChecker) name() string     { return "forward-refs" }
func (c *forwardRefsChecker) severity() string { return "warning" }

func (c *forwardRefsChecker) check(_ string) ([]Violation, error) {
	// TODO: Phase 2
	return []Violation{}, nil
}

// codeAnnotationsChecker validates Go file feature annotations.
type codeAnnotationsChecker struct{}

func newCodeAnnotationsChecker() checker {
	return &codeAnnotationsChecker{}
}

func (c *codeAnnotationsChecker) name() string     { return "code-annotations" }
func (c *codeAnnotationsChecker) severity() string { return "warning" }

func (c *codeAnnotationsChecker) check(_ string) ([]Violation, error) {
	// TODO: Phase 2
	return []Violation{}, nil
}
