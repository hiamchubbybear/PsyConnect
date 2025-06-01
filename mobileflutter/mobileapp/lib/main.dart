import 'package:PsyConnect/core/theme/app_theme.dart';
import 'package:PsyConnect/provider/auth_token_provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/provider/user_profile_provider.dart';
import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/ui/screens/my_home_page.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:provider/provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await dotenv.load(fileName: ".env");
  await Firebase.initializeApp();
  runApp(MultiProvider(
    providers: [
      ChangeNotifierProvider(create: (context) => UserProvider()),
      ChangeNotifierProvider(create: (context) => UserProfileProvider()),
      ChangeNotifierProvider(create: (context) => AuthTokenProvider()),
      ChangeNotifierProvider(create: (_) => ThemeProvider()),
    ],
    child: MyApp(),
  ));
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    final themeProvider = Provider.of<ThemeProvider>(context);
    return MaterialApp(
        debugShowCheckedModeBanner: false,
        theme: ThemesApp.light,
        darkTheme: ThemesApp.dark,
        themeMode: themeProvider.themeMode,
        supportedLocales: const [
          Locale('vi', 'VN'),
          Locale('en', 'US'),
        ],
        locale: const Locale('vi', 'VN'),
        localizationsDelegates: const [
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        home: const MyHomePage(
          title: '',
        ));
  }
}
