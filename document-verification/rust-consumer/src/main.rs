use chrono::Utc;
use kafka::consumer::{Consumer, FetchOffset};
use mongodb::{
    bson::{doc, Document},
    options::ClientOptions,
    sync::Client,
};

use std::str;
use urlencoding::encode;

fn main() {
    let mongo_client = db_init();
    let student_collection: mongodb::sync::Collection<Document> =
        mongo_client.database("University").collection("Students");

    let mut kafka_consumer = queue_init();

    loop {
        for ms in kafka_consumer.poll().unwrap().iter() {

            let applications = ms.messages();
            kafka_consumer.consume_messageset(ms).unwrap();
            kafka_consumer.commit_consumed().unwrap();

            for m in applications {
                println!("{:?}", str::from_utf8(m.value).unwrap());
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
        }
    }
}

fn db_init() -> Client {
    let mut mongo_user = String::new();
    let mut mongo_password = String::new();

    // Take Username and Password from Verifier
    std::io::stdin().read_line(&mut mongo_user).unwrap();
    std::io::stdin().read_line(&mut mongo_password).unwrap();

    // MongoDB connection URI
    let uri = format!(
        "mongodb+srv://{}:{}@hosteldb.e3ayhyn.mongodb.net/?retryWrites=true&w=majority",
        encode(mongo_user.trim()),
        encode(mongo_password.trim())
    );

    // Get a handle to the cluster
    let client_options = ClientOptions::parse(uri).unwrap();
    let client = Client::with_options(client_options).unwrap();

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
        .with_group("verifiers".to_owned())
        .with_fallback_offset(FetchOffset::Earliest)
        .create()
        .unwrap()
}
