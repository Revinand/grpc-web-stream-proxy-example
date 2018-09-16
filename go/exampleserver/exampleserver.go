package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
  "time"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"

	library "github.com/improbable-eng/grpc-web/example/go/_proto/examplecom/library"
)

var (
	grpcEndpoint = "localhost:9090"
)

type bookService struct{}

var books = []*library.Book{
	&library.Book{
		Isbn:   60929871,
		Title:  "Brave New World",
		Author: "Aldous Huxley",
	},
	&library.Book{
		Isbn:   140009728,
		Title:  "Nineteen Eighty-Four",
		Author: "George Orwell",
	},
	&library.Book{
		Isbn:   9780140301694,
		Title:  "Alice's Adventures in Wonderland",
		Author: "Lewis Carroll",
	},
	&library.Book{
		Isbn:   140008381,
		Title:  "Animal Farm",
		Author: "George Orwell",
	},
}

var bookStream = make(chan *library.Book)

func main() {
	flag.Parse()

	grpcServer := grpc.NewServer()
	library.RegisterBookServiceServer(grpcServer, &bookService{})
	grpclog.SetLogger(log.New(os.Stdout, "exampleserver: ", log.LstdFlags))

	//start listener for gRPC server
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatalf("TCP listener error: %s", err.Error())
	}

	//graceful shutdown code:

	// We want to report when each listener is closed.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the grpc listener.
	go func() {
		log.Printf("startup : API server Listening on %s", grpcEndpoint)
		log.Printf("shutdown : API server Listener closed : %v", grpcServer.Serve(lis))
		wg.Done()
	}()

  go func() {
    for {
      bookStream <- books[0]
      time.Sleep(time.Duration(2) * time.Second)
    }
  }()

	// Listen for an interrupt signal from the OS.
	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, os.Interrupt)

	// Wait for a signal to shutdown.
	<-osSignals

	grpcServer.GracefulStop()

	// Wait for the listeners to report it is closed.
	wg.Wait()
	log.Println("main : Completed")

}

func (s *bookService) GetBook(ctx context.Context, bookQuery *library.GetBookRequest) (*library.Book, error) {

	log.Println("GetBook")
	//grpc.SendHeader(ctx, metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-unary"))
	//grpc.SetTrailer(ctx, metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-unary"))

	for _, book := range books {
		if book.Isbn == bookQuery.Isbn {
			return book, nil
		}
	}

	return nil, grpc.Errorf(codes.NotFound, "Book could not be found")
}

func (s *bookService) QueryBooks(bookQuery *library.QueryBooksRequest, stream library.BookService_QueryBooksServer) error {

	log.Println("QueryBooks")
	//stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

  for b := range bookStream {
    if err := stream.Send(b); err != nil {
      log.Printf("error %v", err)
    }
  }

	//stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}
