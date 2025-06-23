plugins {
    java
    id("org.springframework.boot") version "3.1.5"
    id("io.spring.dependency-management") version "1.1.3"
}

group = "com.example"
version = "1.0.0"
java.sourceCompatibility = JavaVersion.VERSION_11

repositories {
    mavenCentral()
}

dependencies {
    // Monitoring Services
    implementation("io.sentry:sentry-spring-boot-starter:6.30.0")
    implementation("com.bugsnag:bugsnag-spring:3.7.0")
    implementation("com.rollbar:rollbar-spring-boot-webmvc:1.10.0")
    implementation("com.datadoghq:datadog-slf4j:1.15.0")
    implementation("com.newrelic.agent.java:newrelic-api:8.5.0")
    implementation("io.micrometer:micrometer-core:1.12.0")

    // Analytics Services
    implementation("com.mixpanel:mixpanel-java:1.5.0")
    implementation("com.amplitude:amplitude-java:1.8.0")
    implementation("com.segment.analytics.java:analytics-java:3.4.0")

    // Payment Services
    implementation("com.stripe:stripe-java:22.30.0")
    implementation("com.paypal:paypal-core:1.7.4")

    // Email Services
    implementation("com.sendgrid:sendgrid-java:4.9.3")
    implementation("net.sargue.mailgun:mailgun-java:1.10.0")
    implementation("com.wildbit.java:postmark-java:1.9.0")

    // Cloud Services
    implementation("software.amazon.awssdk:aws-sdk-java:2.20.162")
    implementation("com.google.cloud:google-api-services-storage:v1-rev20230301-2.0.0")
    implementation("com.microsoft.azure:azure-sdk-bom:1.2.18")

    // Standard Spring Boot dependencies
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    implementation("org.springframework.boot:spring-boot-starter-security")

    // Testing
    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("junit:junit:4.13.2")
    testImplementation("org.mockito:mockito-core:5.5.0")
}

tasks.named<Test>("test") {
    useJUnitPlatform()
}