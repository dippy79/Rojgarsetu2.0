# Flutter Screens Phase F - 7 Missing Screens
Status: ✅ PLAN APPROVED | 🔄 IMPLEMENTING

## Steps (sequential):
- [ ✅ ] 1. Check/add url_launcher to pubspec.yaml (already present)
- [ ✅ ] 2. Edit app.dart - Add GoRoutes
- [ ✅ ] 3. Create candidate_profile_screen.dart
- [ ✅ ] 4. Create company_profile_screen.dart  
- [ ✅ ] 5. Create my_applications_screen.dart
- [ ✅ ] 6. Create job_applications_screen.dart
- [ ✅ ] 7. Create gov_jobs_screen.dart
- [ ✅ ] 8. Create courses_screen.dart
- [ ✅ ] 9. Create videos_screen.dart
- [ 🔄 ] 10. flutter pub get & flutter analyze (0 errors)
- [ ] 11. Update this TODO.md to ✅ COMPLETE
- [x] HomeScreen already imports gov_jobs_screen.dart etc. as tabs

## Notes:
- BLoC for lists (GovJobsBloc etc. assumed complete), direct ApiService for profiles/applications
- url_launcher for external links (gov jobs notificationUrl, courses/videos urls)
- StatusBadge colors match Application.statusColor
- Loading/error/empty states everywhere
- SmartRefresher for PTR + pagination

