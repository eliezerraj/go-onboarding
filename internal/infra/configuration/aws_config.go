package configuration

import(
	"os"

	"github.com/joho/godotenv"
	"github.com/go-onboarding/internal/core/model"
)

// About get AWS service env ver
func GetAwsServiceEnv() model.AwsService {
	childLogger.Info().Str("func","GetAwsServiceEnv").Send()

	err := godotenv.Load(".env")
	if err != nil {
		childLogger.Info().Err(err).Send()
	}
	
	var awsService	model.AwsService

	if os.Getenv("AWS_REGION") !=  "" {
		awsService.AwsRegion = os.Getenv("AWS_REGION")
	}

	if os.Getenv("BUCKET_NAME") !=  "" {
		awsService.BucketName = os.Getenv("BUCKET_NAME")
	}

	if os.Getenv("FILE_PATH") !=  "" {
		awsService.FilePath = os.Getenv("FILE_PATH")
	}

	return awsService
}