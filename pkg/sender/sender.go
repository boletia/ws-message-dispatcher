package sender

type sender struct {
	awsRegion      string
	lambdaFunction string
}

// New Creates new sender instance
func New(region, funcName string) sender {
	return sender{
		awsRegion:      region,
		lambdaFunction: funcName,
	}
}
