# Flutter Production Release Build - RojgarSetu

## APK (direct distribution)

```bash
cd mobile_app_flutter

# 1. Update baseUrl for production
API_BASE_URL=https://api.yourdomain.com flutter pub get

# 2. Build APK
flutter build apk --release

# Output: build/app/outputs/flutter-apk/app-release.apk ~20MB
```

## App Bundle (Play Store)

```bash
flutter build appbundle --release

# Output: build/app/outputs/bundle/release/app-release.aab
```

## Signing Setup (REQUIRED for Play Store)

1. Generate keystore:
```bash
keytool -genkey -v -keystore android/rojgarsetu.jks -keyalg RSA -keysize 2048 -validity 10000 -alias rojgarsetu
```

2. Create `android/key.properties`:
```
storePassword=your_store_password
keyPassword=your_key_password
keyAlias=rojgarsetu
storeFile=../android/rojgarsetu.jks
```

3. Update `android/app/build.gradle` (add before android {):
```
def keystoreProperties = new Properties()
def keystorePropertiesFile = rootProject.file('key.properties')
if (keystorePropertiesFile.exists()) {
    keystoreProperties.load(new FileInputStream(keystorePropertiesFile))
}

android {
    ...
    signingConfigs {
        release {
            keyAlias keystoreProperties['keyAlias']
            keyPassword keystoreProperties['keyPassword']
            storeFile keystoreProperties['storeFile'] ? file(keystoreProperties['storeFile']) : null
            storePassword keystoreProperties['storePassword']
        }
    }
    buildTypes {
        release {
            signingConfig signingConfigs.release
            ...
        }
    }
}
```

## CI/CD APK (GitHub Actions)
- Push to main → auto builds APK (unsigned)
- Download from Actions artifacts

## Notes
- baseUrl from env API_BASE_URL
- Signing optional for APK sideloading
- Test on physical device before Play Store
- ProGuard enabled by default (--no-shrink false)

