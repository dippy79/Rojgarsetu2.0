package com.rojgarsetu.auth;

import java.nio.charset.StandardCharsets;
import java.util.Date;

import javax.crypto.SecretKey;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
// SignatureAlgorithm is removed because 0.12.x auto-determines it from the Key, or you can specify the Jwts.SIG algorithm directly.
import io.jsonwebtoken.security.Keys;

/**
 * JWT generation/validation for the auth service.
 *
 * IMPORTANT (interoperability): tokens produced here are designed to be
 * accepted by the Go backend (backend_go) AuthMiddleware, so the two services
 * share:
 *   - the same signing secret  : JWT_SECRET (env), required, >= 32 bytes
 *   - the same algorithm       : HS256
 *   - the same custom claims   : user_id, email, role
 *   - the same issuer/audience : JWT_ISSUER (default "rojgarsetu-backend"),
 *                                JWT_AUDIENCE (default "rojgarsetu-api")
 *
 * If you change any of these defaults here, update backend_go/config/config.go
 * and backend_go/internal/middleware/auth.go to match, otherwise the Go backend
 * will reject tokens issued by this service (issuer/audience mismatch).
 */
@Service
public class JwtService {

    private final SecretKey signingKey;
    private final String issuer;
    private final String audience;
    private final long expirationMs;

    public JwtService(
            @Value("${jwt.secret}") String jwtSecret,
            @Value("${jwt.issuer:rojgarsetu-backend}") String issuer,
            @Value("${jwt.audience:rojgarsetu-api}") String audience,
            @Value("${jwt.expiry-seconds:86400}") long expirySeconds) {
        if (jwtSecret == null || jwtSecret.isEmpty()) {
            throw new IllegalStateException("JWT_SECRET environment variable is required (set jwt.secret)");
        }
        if (jwtSecret.getBytes(StandardCharsets.UTF_8).length < 32) {
            throw new IllegalStateException("JWT_SECRET must be at least 32 bytes long for security");
        }
        this.signingKey = Keys.hmacShaKeyFor(jwtSecret.getBytes(StandardCharsets.UTF_8));
        this.issuer = issuer;
        this.audience = audience;
        this.expirationMs = expirySeconds * 1000L;
    }

    /**
     * Generate a signed HS256 JWT with standard + custom claims that match the
     * Go backend's Claims struct (user_id, email, role).
     */
    public String generateToken(String userId, String email, String role) {
        Date now = new Date();
        Date expiry = new Date(now.getTime() + expirationMs);
        return Jwts.builder()
                .subject(userId)                 // Updated: setSubject -> subject
                .claim("user_id", userId)
                .claim("email", email)
                .claim("role", role)
                .issuer(issuer)                  // Updated: setIssuer -> issuer
                .audience().add(audience).and()  // Updated: setAudience -> audience().add(...).and()
                .issuedAt(now)                   // Updated: setIssuedAt -> issuedAt
                .expiration(expiry)              // Updated: setExpiration -> expiration
                .signWith(signingKey)            // Updated: signWith(Key, Alg) -> signWith(Key)
                .compact();
    }

    /**
     * Validate and parse a token. Throws JwtException on invalid/expired tokens.
     */
    public Claims validateToken(String token) {
        return Jwts.parser()                     // Updated: parserBuilder() -> parser()
                .verifyWith(signingKey)          // Updated: setSigningKey() -> verifyWith()
                .build()
                .parseSignedClaims(token)        // Updated: parseClaimsJws() -> parseSignedClaims()
                .getPayload();                   // Updated: getBody() -> getPayload()
    }
}