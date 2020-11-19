package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/api/option"
)

var db = dbInfo{"root", "1463", "localhost:3306", "mysql", "running_plus"}

type dbInfo struct {
	user     string
	pwd      string
	url      string
	engine   string
	database string
}

/*
type token_view struct {
	view_token string
}
*/

func main() {
	http.HandleFunc("/", index_page)
	http.HandleFunc("/push/", push_page)

	log.Println("Listening on : 9090...")
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServer : ", err)
	} else {
		fmt.Println("ListenAndServer Started! -> Port(9000)")
	}
}

func index_page(w http.ResponseWriter, r *http.Request) {
	indexTemplate, _ := template.ParseFiles("index.html")
	indexTemplate.Execute(w, nil)
}

func push_page(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "POST" {
		restitle := r.FormValue("title")
		rescontent := r.FormValue("content")
		if (restitle != "") && (rescontent != "") {
			fmt.Println("OK")

			var token_select = "SELECT Token FROM users WHERE id=3;"
			var db_token = SelectQuery(db, token_select)

			opt := option.WithCredentialsFile("running-plus-9c5ab-firebase-adminsdk-ndi06-ee29560d12.json")
			app, _ := firebase.NewApp(context.Background(), nil, opt)
			sendToToken(app, restitle, rescontent, db_token)

			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			fmt.Println("Access Denied !")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func SelectQuery(db dbInfo, query string) string {
	var db_token string
	dataSource := db.user + ":" + db.pwd + "@tcp(" + db.url + ")/" + db.database
	conn, err := sql.Open(db.engine, dataSource)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := conn.Query(query)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&db_token)
		if err != nil {
			log.Fatal(err)
		}
	}
	return db_token
}

func initializeAppDefault() *firebase.App {
	// [START initialize_app_default_golang]
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_app_default_golang]

	return app
}

func sendToTopic(ctx context.Context, client *messaging.Client, title string, content string, dbtoken string) {
	// [START send_to_topic_golang]
	// The topic name can be optionally prefixed with "/topics/".

	//token := "fY_9WFgkz7g:APA91bGYAslhXrRU_u9V0rNkXEu0QEfNgU2RdoNZ6c8uDCtolrJLbRhDu9iPV-TDfUv78x-HiJ6L203WD9R8sys_xS4oCumgs9MEnvyTxAhy-o1LA2d5jFkSxvJd38mIVEZk5QZ4Ev84"

	//token := dbtoken
	//topic := "running_plus_fcm"
	topic := "newcomers_fcm"

	// See documentation on defining a message payload.
	messages := []*messaging.Message{
		{

			Notification: &messaging.Notification{
				Title: title,
				Body:  content,
			},

			Data: map[string]string{
				"title": title,
				"body":  content,
			},
			//Token: token,
			Topic: topic,
		},
	}

	// Send a message to the devices subscribed to the provided topic.
	br, err := client.SendAll(context.Background(), messages)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", br.SuccessCount)
	// [END send_to_topic_golang]
}

func initializeAppWithRefreshToken() *firebase.App {
	// [START initialize_app_refresh_token_golang]
	opt := option.WithCredentialsFile("path/to/refreshToken.json")
	config := &firebase.Config{ProjectID: "running-plus-9c5ab"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_app_refresh_token_golang]

	return app
}

func sendToToken(app *firebase.App, title string, content string, dbtoken string) {
	// [START send_to_token_golang]
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	//registrationToken := "fY_9WFgkz7g:APA91bGYAslhXrRU_u9V0rNkXEu0QEfNgU2RdoNZ6c8uDCtolrJLbRhDu9iPV-TDfUv78x-HiJ6L203WD9R8sys_xS4oCumgs9MEnvyTxAhy-o1LA2d5jFkSxvJd38mIVEZk5QZ4Ev84"
	registrationToken := dbtoken

	//topic := "running_plus_fcm"
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"title": title,
			"body":  content,
			//"score": title,
			//"time":  content,
		},
		Token: registrationToken,
		//Topic: topic,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	// [END send_to_token_golang]
}

/*
func sendAll(ctx context.Context, client *messaging.Client, title string, content string) {
	// This registration token comes from the client FCM SDKs.
	registrationToken := "ccm06rM6pHw:APA91bFD8b9AB_IJsKve_uTEU9JE8bbDK0c6MCmtPtwCaKB1aE2573fFWAPWCw36HDqBFF3w2XQPxbZfV_nB-TCcbb53MSbIBryDduWEYJwWomxrbjTSfTkxWH7JKWF4HdHZcJh4AtOr"

	// [START send_all]
	// Create a list containing up to 100 messages.
	messages := []*messaging.Message{
		{
			Notification: &messaging.Notification{
				Title: title,
				Body:  content,
			},
			Token: registrationToken,
		},
	}

	br, err := client.SendAll(context.Background(), messages)
	if err != nil {
		log.Fatalln(err)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_all]
}
*/
