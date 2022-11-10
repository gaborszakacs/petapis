package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	// https://buf.build/gaborszakacs-bitrise/pet/docs/main:pet.v1
	// https://buf.build/bufbuild/templates/connect-go/v10
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	petv1 "go.buf.build/bufbuild/connect-go/gaborszakacs-bitrise/pet/pet/v1"
	"go.buf.build/bufbuild/connect-go/gaborszakacs-bitrise/pet/pet/v1/petv1connect"
)

type PetServer struct{}

func (s *PetServer) GetPet(
	ctx context.Context,
	req *connect.Request[petv1.GetPetRequest],
) (*connect.Response[petv1.GetPetResponse], error) {
	log.Println("Request headers: ", req.Header())
	id := req.Msg.GetPetId()

	log.Println("Got a request to get pet: ", id)

	dogNames := []string{"Fido", "Spot", "Buddy", "Barkley", "Barkley the Third", "Coco", "Max", "Duke", "Bella", "Ruby", "Bandit", "Stella"}
	name := dogNames[rand.Intn(len(dogNames))]

	res := connect.NewResponse(&petv1.GetPetResponse{
		Pet: &petv1.Pet{Name: name, PetId: id, PetType: petv1.PetType_PET_TYPE_DOG},
	})
	res.Header().Set("Pet-Version", "v1")
	return res, nil
}

func (s *PetServer) GetPetStream(
	ctx context.Context,
	req *connect.Request[petv1.GetPetStreamRequest],
	stream *connect.ServerStream[petv1.GetPetStreamResponse],
) error {
	for i := 0; i < 5; i++ {
		dogNames := []string{"Fido", "Spot", "Buddy", "Barkley", "Barkley the Third", "Coco", "Max", "Duke", "Bella", "Ruby", "Bandit", "Stella"}
		random := rand.Intn(len(dogNames))
		name := dogNames[random]

		res := &petv1.GetPetStreamResponse{
			Pet: &petv1.Pet{Name: name, PetId: fmt.Sprintf("%d", random), PetType: petv1.PetType_PET_TYPE_DOG},
		}

		fmt.Printf("\n%d) Sending pet, but sleeping for %d sec first", i+1, random)
		time.Sleep(time.Duration(random) * time.Second)

		stream.Send(res)
	}

	return nil
}

func main() {
	server := &PetServer{}
	mux := http.NewServeMux()
	path, handler := petv1connect.NewPetStoreServiceHandler(server)

	c := cors.AllowAll()
	handler = c.Handler(handler)

	// Offer the gRPC reflection service
	mux.Handle(path, handler)
	reflector := grpcreflect.NewStaticReflector(
		"pet.v1.PetStoreService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

// dummy implementation
func (s *PetServer) PutPet(
	ctx context.Context,
	req *connect.Request[petv1.PutPetRequest],
) (*connect.Response[petv1.PutPetResponse], error) {
	return nil, nil
}

// dummy implementation
func (s *PetServer) DeletePet(
	ctx context.Context,
	req *connect.Request[petv1.DeletePetRequest],
) (*connect.Response[petv1.DeletePetResponse], error) {
	return nil, nil
}
