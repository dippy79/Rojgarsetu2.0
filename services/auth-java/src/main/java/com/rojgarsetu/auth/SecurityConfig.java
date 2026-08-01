package com.rojgarsetu.auth;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.provisioning.InMemoryUserDetailsManager;
import org.springframework.security.web.SecurityFilterChain;

/**
 * Spring Security configuration for the auth service.
 *
 * This is a stateless REST API that authenticates via its own /auth/login and
 * /auth/register endpoints (BCrypt + JWT-style tokens). The default Spring
 * Security configuration (CSRF protection, session-based auth, generated
 * password) is wrong for this model, so we replace it with a custom
 * SecurityFilterChain:
 *
 *  - CSRF disabled       – CSRF is for browser form sessions, not stateless JWT REST APIs
 *  - STATELESS sessions  – no server-side session state
 *  - permitAll           – auth endpoints and the actuator health probe
 *  - anyRequest          – everything else requires authentication
 */
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
            // CSRF protection is for browser form sessions, not stateless JWT REST APIs.
            .csrf(csrf -> csrf.disable())
            // No server-side sessions – every request is authenticated independently.
            .sessionManagement(session ->
                session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/auth/**", "/auth/**", "/actuator/health").permitAll()
                .anyRequest().authenticated());
        return http.build();
    }

    /**
     * BCrypt password encoder shared by the registration and login flows.
     * (AuthController currently instantiates one inline; this bean makes the
     * encoder injectable and avoids duplicating configuration.)
     */
    @Bean
    public BCryptPasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    /**
     * Provide a UserDetailsService so Spring Boot does NOT auto-configure its
     * default in-memory user and suppress the "Using generated security
     * password: ..." startup warning. This service does not rely on Spring
     * Security's authentication manager – login is handled by the JPA-backed
     * /auth/login endpoint – so an empty user store is intentional.
     */
    @Bean
    public UserDetailsService userDetailsService() {
        return new InMemoryUserDetailsManager();
    }
}

