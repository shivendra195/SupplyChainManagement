package env

const (
	// BRANCH is the code's current git branch.
	BRANCH string = "BRANCH"

	// DEV is the dev branch name
	DEV string = "dev"

	// MAIN is the prod branch name
	MAIN string = "main"
)

var kubernetesServiceHost string
var currentBranch string

func IsMain() bool {
	return IsBranch(MAIN)
}

func IsDev() bool {
	return IsBranch(DEV)
}

func InKubeCluster() bool {
	return kubernetesServiceHost != ""
}

// IsBranch checks if the current code is part of the specified branch (name).
func IsBranch(name string) bool {
	return name == currentBranch
}
