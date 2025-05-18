package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func GrpcProvider(logger sflogger.Logger) {
	//err := grpco.RegisterConnection(
	//	// Define the connection details and pass services directly
	//	grpco.WithConnectionDetails(
	//		"story-service",
	//		"story-service.example.com:50051",
	//		grpco.WithInsecure(),
	//		map[string]grpco.ServiceDefinition{
	//			"story": {
	//				ClientConstructor: svc.NewStoryServiceClient,
	//				Methods: map[string]string{
	//					"List": "/service.StoryService/List",
	//					"Get":  "/service.StoryService/Get",
	//				},
	//			},
	//		}, // Pass services as an argument
	//	),
	//	grpco.WithLogger(logger),
	//)
	//if err != nil {
	//	logger.Errorf("Failed to register gRPC connection : %s", "error", err.Error())
	//}
}
