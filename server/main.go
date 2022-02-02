package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"mime"
	"net"
	"net/http"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	"github.com/gaborszakacs/petapis/docs"

	// petv1 "go.buf.build/grpc/go/gaborszakacs/petapis/pet/v1"
	petv1 "github.com/gaborszakacs/petapis/gen/proto/go/pet/v1"

	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	go runOpenAPIFileServer()
	listenOn := "127.0.0.1:8080"
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", listenOn, err)
	}

	server := grpc.NewServer()
	petv1.RegisterPetStoreServiceServer(server, &petStoreServiceServer{})
	log.Println("Listening on", listenOn)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

// petStoreServiceServer implements the PetStoreService API.
type petStoreServiceServer struct {
	petv1.UnimplementedPetStoreServiceServer
}

// PutPet adds the pet associated with the given request into the PetStore.
func (s *petStoreServiceServer) PutPet(ctx context.Context, req *petv1.PutPetRequest) (*petv1.PutPetResponse, error) {
	name := req.GetName()
	petType := req.GetPetType()
	log.Println("Got a request to create a", petType, "named", name)

	return &petv1.PutPetResponse{}, nil
}

func (s *petStoreServiceServer) GetPet(ctx context.Context, req *petv1.GetPetRequest) (*petv1.GetPetResponse, error) {
	id := req.GetPetId()
	log.Println("Got a request to get ", id)

	dogNames := []string{"Fido", "Spot", "Buddy", "Barkley", "Barkley the Third"}
	name := dogNames[rand.Intn(len(dogNames))]

	return &petv1.GetPetResponse{
		Pet: &petv1.Pet{Name: name, PetId: id, PetType: petv1.PetType_PET_TYPE_DOG},
	}, nil
}

func runOpenAPIFileServer() {
	fileHandler := getOpenAPIHandler()
	mux := http.NewServeMux()
	mux.Handle("/docs/", http.StripPrefix("/docs/", fileHandler))
	port := "8081"
	addr := "127.0.0.1:" + port
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(docs.OpenAPI, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}
