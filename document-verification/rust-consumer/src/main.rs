use chrono::Utc;
use kafka::consumer::{Consumer, FetchOffset};
use mongodb::{
    bson::{doc, Document},
    sync::Client,
};
use std::env;
use std::str;

fn main() {
    let mongo_client = db_init();
    let student_collection: mongodb::sync::Collection<Document> =
        mongo_client.database("University").collection("Students");

    let mut kafka_consumer = queue_init();

    loop {
        for ms in kafka_consumer.poll().unwrap().iter() {
            for m in ms.messages() {
                let (user, form_id) = str::from_utf8(m.value).unwrap().split_once(':').unwrap();
                println!("{:?}", form_id);
                student_collection
                    .update_one(
                        doc! {
                            "forms": {
                                "$exists": form_id,
                            },
                            "name": user,
                        },
                        doc! {
                            "$set": {
                                format!("forms.{form_id}.status"): "VERIFIED",
                                format!("forms.{form_id}.verification"): Utc::now().to_rfc3339(),
                            },
                        },
                        None,
                    )
                    .unwrap();
            }
            match kafka_consumer.consume_messageset(ms) {
                Ok(number) => number,
                Err(e) => println!("{:?}", e),
            };
        }
        match kafka_consumer.commit_consumed() {
            Ok(number) => number,
            Err(e) => println!("{:?}", e),
        };
    }
}

fn db_init() -> Client {
    // Path to certificate
    let mongo_certificate = env::var("MONGODB_CERTIFICATE").unwrap();

    // MongoDB connection URI
    let uri = "mongodb+srv://hosteldb.e3ayhyn.mongodb.net/?\
         retryWrites=true&w=majority \
         &authSource=%24external&authMechanism=MONGODB-X509&tlsCertificateKeyFile="
        .to_owned()
        + &mongo_certificate;

    // Get a handle to the cluster
    let client = Client::with_uri_str(uri).unwrap();

    // Ping the server to see if you can connect to the cluster
    client
        .database("University")
        .run_command(doc! {"ping": 1}, None)
        .unwrap();
    println!("Connected successfully.");

    client
}

fn queue_init() -> Consumer {
    Consumer::from_hosts(vec!["172.27.0.6:9092".to_owned()])
        .with_topic("test".to_owned())
        .with_fallback_offset(FetchOffset::Earliest)
        .create()
        .unwrap()
}
