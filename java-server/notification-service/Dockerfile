FROM eclipse-temurin:17-jdk-jammy
WORKDIR /app
COPY pom.xml /app/
COPY src /app/src
RUN mvn -B package -DskipTests
EXPOSE 8080
CMD ["java", "-cp", "target/notification-service-0.1.0.jar", "com.example.notification.App"]
