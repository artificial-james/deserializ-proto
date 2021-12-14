package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/artificialinc/alab-core/common-go/stream/redis"
	pb "github.com/artificialinc/artificial-protos/go/artificial/api/alab/scheduler"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "invalid arguments\n")
		os.Exit(1)
	}
	namespace := os.Args[1]
	orgID := os.Args[2]
	jobID := os.Args[3]

	stream := redis.NewStream(nil)
	payloads, _, err := stream.GetLatest(
		context.Background(),
		fmt.Sprintf("%s:org.%s.job.%s.program.0", namespace, orgID, jobID),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get latest\n")
		os.Exit(1)
	}

	var b string
	for _, payload := range payloads {
		for _, v := range payload {
			b = v
			break
		}
	}
	a := &pb.Graph{}

	err = proto.Unmarshal([]byte(b), a)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot unmarshal to json: %v\n", err)
		os.Exit(1)
	}

	j := protojson.Format(a)

	fmt.Println(j)

	fmt.Fprintf(os.Stderr, "binary size: %d\n", len(b))
	fmt.Fprintf(os.Stderr, "number of nodes: %d\n", len(a.Nodes))
	fmt.Fprintf(os.Stderr, "number of edges: %d\n", len(a.Edges))
}
