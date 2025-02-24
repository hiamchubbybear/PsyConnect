import 'package:PsyConnect/page/login_page.dart';
import 'package:PsyConnect/provider/auth_token_provider.dart';
import 'package:PsyConnect/provider/user_profile_provider.dart';
import 'package:PsyConnect/provider/user_provider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:provider/provider.dart';

void main() async {
  await dotenv.load(fileName: ".env");
  runApp(MultiProvider(
    providers: [
      ChangeNotifierProvider(create: (context) => UserProvider()),
      ChangeNotifierProvider(create: (context) => UserProfileProvider()),
      ChangeNotifierProvider(create: (context) => AuthTokenProvider()),
    ],
    child: const MyApp(),
  ));
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(
              seedColor: const Color.fromARGB(255, 125, 170, 249)),
          useMaterial3: true,
        ),
        supportedLocales: const [
          Locale('vi', 'VN'),
          Locale('en', 'US'),
        ],
        locale: const Locale('vi', 'VN'),
        localizationsDelegates: [
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        home:
            // const VerifiedPage()
            // const MyHomePage(title:  "PsyConnect",),
            const LoginPage());
  }
}
