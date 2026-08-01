import 'package:flutter/material.dart';
class AppTheme {
   static ThemeData get lightTheme => ThemeData(
       useMaterial3: true,
        colorScheme: ColorScheme.fromSeed(
           seedColor: const Color(0xFF1E88E5), 
            brightness: Brightness.light,
            primary: const Color(0xFF1E88E5),
             secondary: const Color(0xFFFF8F00),
               ),
             scaffoldBackgroundColor: Colors.grey[50],
            appBarTheme: const AppBarTheme(
            elevation: 0,
            centerTitle: true,
             titleTextStyle: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.white,
               ),
              ),
              elevatedButtonTheme: ElevatedButtonThemeData(
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                     ), ),
                      ),
                      inputDecorationTheme: InputDecorationTheme(
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),),
                          filled: true,
                          fillColor: Colors.white,
                          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                          ),
                          );
                          static ThemeData get darkTheme => ThemeData(
                            useMaterial3: true,
                            colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF1E88E5),
                            brightness: Brightness.dark,
                            ),
                            scaffoldBackgroundColor: Colors.grey[900],
                            );
                            }
