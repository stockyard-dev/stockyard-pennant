package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-pennant/internal/server";"github.com/stockyard-dev/stockyard-pennant/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="10090"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./pennant-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("pennant: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Pennant — loyalty and points system\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("pennant: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
