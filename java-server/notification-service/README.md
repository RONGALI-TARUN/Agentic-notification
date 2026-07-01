Java production example

Build and run (requires JDK 11+):

```bash
cd java-server/production
javac -d out $(find src/main/java -name "*.java")
java -cp out com.example.notification.App
```

POST to `http://localhost:8080/notify` with the JSON format used earlier.
