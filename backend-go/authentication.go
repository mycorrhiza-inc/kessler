package main

import (
	"os"
)

// FIXME : HIGHLY INSECURE, GET THE HMAC SECRET FROM SUPABASE AND THROW IT IN HERE AS AN NEV VARAIBLE.
var SupabaseSecret = os.Getenv("SUPABASE_ANON_KEY")

// func makeTokenValidator(dbtx_val dbstore.DBTX) func(r *http.Request) UserValidation {
// 	return_func := func(r *http.Request) UserValidation {
// 		token := r.Header.Get("Authorization")
// 		if token == "" {
// 			// Check if the supabase cookie starting with sb-kpvkpczxcclxslabfzeu-auth-token, exists if it does do the decoding and set token equal to Bearer <jwt_token>, otherwise return an anonomous ser
// 			cookie, err := r.Cookie("sb-kpvkpczxcclxslabfzeu-auth-token")
// 			if err != nil {
// 				return UserValidation{
// 					validated: true,
// 					userID:    "anonomous",
// 				}
// 			}
// 			// Strip prefix and decode Base64 part.
// 			if !strings.HasPrefix(cookie.Value, "base64-") {
// 				// Json will catch if invalid
// 				fmt.Println("Cookie is not base64 decodable.")
// 			}
// 			encodedData := strings.TrimSpace(strings.TrimPrefix(cookie.Value, "base64-"))
// 			decodedData, err := base64.URLEncoding.DecodeString(encodedData)
// 			if err != nil {
// 				// Json will catch if invalid
// 				fmt.Printf("Error decoding base64 %v\n", err)
// 			}
// 			stringData := string(decodedData)
// 			// var tokenData AccessTokenData
// 			// TODO : Fix horrible moneky wrench solution for decoding this with something not shit
// 			// err = json.Unmarshal([]byte(string(decodedData)), &tokenData)
// 			// if err != nil {
// 			// 	fmt.Printf("Error unmarshalling %v\n", err)
// 			// 	return UserValidation{validated: false}
// 			// }
// 			// token = fmt.Sprintf("Bearer %s", tokenData.AccessToken)
// 			stringDataStripped := strings.TrimPrefix(stringData, `{"access_token":"`)
// 			hopefullyToken := strings.Split(stringDataStripped, `"`)[0]
// 			_ = hopefullyToken // here to prevent the compiler from complaining
// 			// token = fmt.Sprintf("Bearer %s", hopefullyToken)
// 			// fmt.Println(token)
//
// 		}
// 		// Check for "Bearer " prefix in the authorization header (expected format)
// 		if !strings.HasPrefix(token, "Bearer ") {
// 			return UserValidation{
// 				validated: true,
// 				userID:    "anonomous",
// 			}
// 		}
// 		// Validation four our scrapers to add data to the system
// 		if strings.HasPrefix(token, "Bearer thaum_") {
// 			return UserValidation{userID: "thaumaturgy", validated: true}
// 			// TODO: Add a check so that authentication only succeeds if it comes from a tailscale IP.
// 			// 		q := *routing.DBQueriesFromRequest(r)

// 			// const trim = len("Bearer thaum_")
// 			// // Replacing this with PBKDF2 or something would be more secure, but it should matter since every API key can be gaurenteed to have at least 128/256 bits of strength.
// 			// hash := blake2b.Sum256([]byte(token[trim:]))
// 			// encodedHash := base64.URLEncoding.EncodeToString(hash[:])
// 			// fmt.Println("Checking Database for Hashed API Key:", encodedHash)
// 			// ctx := r.Context()
// 			// result, err := q.CheckIfThaumaturgyAPIKeyExists(ctx, encodedHash)
// 			// if result.KeyBlake3Hash == encodedHash && err != nil {
// 			// 	return UserValidation{userID: "thaumaturgy", validated: true}
// 			// }
// 			// return UserValidation{validated: false}
// 		}
//
// 		tokenString := strings.TrimPrefix(token, "Bearer ")
//
// 		// Parse the JWT token
// 		keyFunc := func(token *jwt.Token) (interface{}, error) {
// 			// Validate the algorithm
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			// Return the secret for signature verification
// 			jwtSecret := []byte(SupabaseSecret)
// 			return jwtSecret, nil
// 		}
// 		parsedToken, err := jwt.Parse(tokenString, keyFunc)
// 		// fmt.Println(parsedToken)
// 		if err != nil {
// 			fmt.Printf("Encountered error with token validation since that functionality hasnt been implemented yet, and the backend assumes every HMAC signature is valid, this is probably good to fix if we dont want to get royally screwed %v", err)
// 			// FIXME : HIGHLY INSECURE, GET THE HMAC SECRET FROM SUPABASE AND THROW IT IN HERE AS AN NEV VARAIBLE.
// 			// return UserValidation{validated: false}
// 		}
//
// 		// FIXME : HIGHLY INSECURE, GET THE HMAC SECRET FROM SUPABASE AND THROW IT IN HERE AS AN NEV VARAIBLE.
// 		claims, ok := parsedToken.Claims.(jwt.MapClaims)
//
// 		fmt.Println(claims)
// 		ok = ok || !ok
//
// 		// if ok && parsedToken.Valid {
// 		if ok {
// 			userID := claims["sub"] // JWT 'sub' - typically the user ID
// 			// Perform additional checks if necessary
// 			return UserValidation{userID: userID.(string), validated: true}
// 		}
//
// 		return UserValidation{validated: false}
// 	}
//
// 	return return_func
// }

// func makeAuthMiddleware(dbtx_val dbstore.DBTX) func(http.Handler) http.Handler {
// 	tokenValidator := makeTokenValidator(dbtx_val)
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			userInfo := tokenValidator(r)
// 			if userInfo.validated {
// 				r.Header.Set("Authorization", fmt.Sprintf("Authenticated %s", userInfo.userID))
// 				// fmt.Printf("Authenticated Request for user %v\n", userInfo.userID)
// 				next.ServeHTTP(w, r)
//
// 			} else {
// 				fmt.Println("Auth Failed, for ip address", r.RemoteAddr)
// 				http.Error(w, "Authentication failed", http.StatusUnauthorized)
// 			}
// 		})
// 	}
// }
