package generator

import "embed"

//go:embed assets
var embeddedAssets embed.FS

var templateData = map[string]string{
	"main.dart": `// uncomment lines after configuring you're project with you're firebase console project. 
import 'core/constants/routes.dart';
import 'features/auth/presentation/bloc/auth_bloc.dart';
import 'features/auth/presentation/bloc/auth_event.dart';
import 'features/auth/presentation/pages/auth_wrapper.dart';
// import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:shadcn_ui/shadcn_ui.dart';
// import 'firebase_options.dart';

import 'features/auth/data/repositories/auth_repository_impl.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:firebase_auth/firebase_auth.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  // await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);

  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  final googleSignIn = GoogleSignIn();
  late final authRepository = AuthRepositoryImpl(
    FirebaseAuth.instance,
    googleSignIn,
  );

  final GlobalKey<ScaffoldMessengerState> _scaffoldMessengerKey =
      GlobalKey<ScaffoldMessengerState>();

  MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider(
          create:
              (_) =>
                  AuthBloc(authRepository)
                    ..add(const AuthEvent.checkAuthStatus()),
        ),
      ],

      child: MaterialApp(
        scaffoldMessengerKey: _scaffoldMessengerKey,
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          colorScheme: const ColorScheme.light(),
          appBarTheme: const AppBarTheme(
            backgroundColor: Colors.transparent,
            elevation: 0,
          ),
        ),
        home: Builder(
          builder: (context) {
            return ShadApp(
              debugShowCheckedModeBanner: false,
              darkTheme: ShadThemeData(
                brightness: Brightness.dark,
                colorScheme: const ShadSlateColorScheme.dark(),
              ),
              initialRoute: AppRoutes.initial,
              onGenerateRoute: AppRoutes.onGenerateRoute,
              home: const AuthWrapper(),
            );
          },
        ),
      ),
    );
  }
}
`,

	"widget_test.dart": `// This is a basic Flutter widget test.
// To perform an interaction with a widget in your test, use the WidgetTester
// utility in the flutter_test package. For example, you can send tap and scroll
// gestures. You can also use WidgetTester to find child widgets in the widget
// tree, read text, and verify that the values of widget properties are correct.

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import '../lib/main.dart';

void main() {
  testWidgets('Counter increments smoke test', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(MyApp());

    // Verify that our counter starts at 0.
    expect(find.text('0'), findsOneWidget);
    expect(find.text('1'), findsNothing);

    // Tap the '+' icon and trigger a frame.
    await tester.tap(find.byIcon(Icons.add));
    await tester.pump();

    // Verify that our counter has incremented.
    expect(find.text('0'), findsNothing);
    expect(find.text('1'), findsOneWidget);
  });
}
`,

	"routes.dart": `import 'package:flutter/material.dart';
import '../../features/auth/presentation/pages/signin_page.dart';
import '../../features/auth/presentation/pages/signup_page.dart';
import '../../features/home/presentation/pages/home_page.dart';

class AppRoutes {
  static const String initial = '/';
  static const String signIn = '/sign-in';
  static const String signUp = '/sign-up';
  static const String home = '/home';

  static Route<dynamic> onGenerateRoute(RouteSettings settings) {
    switch (settings.name) {
      case initial:
      case signIn:
        return MaterialPageRoute(builder: (_) => const SignInScreen());
      case signUp:
        return MaterialPageRoute(builder: (_) => const SignUpScreen());
      case home:
        return MaterialPageRoute(builder: (_) => const HomeScreen());
      default:
        return MaterialPageRoute(
          builder:
              (_) => Scaffold(
                appBar: AppBar(title: const Text('Not Found')),
                body: const Center(child: Text('Page not found')),
              ),
        );
    }
  }

  static Future<dynamic> navigateTo(BuildContext context, String routeName) {
    return Navigator.of(context).pushNamed(routeName);
  }

  static Future<dynamic> navigateToAndRemoveUntil(
    BuildContext context,
    String routeName,
  ) {
    return Navigator.of(
      context,
    ).pushNamedAndRemoveUntil(routeName, (Route<dynamic> route) => false);
  }

  static Future<dynamic> navigateToAndReplace(
    BuildContext context,
    String routeName,
  ) {
    return Navigator.of(context).pushReplacementNamed(routeName);
  }

  static void goBack(BuildContext context) {
    Navigator.of(context).pop();
  }
}
`,

	"user_model.dart": `import '../../domain/entities/user_entity.dart';

class UserModel {
  final UserEntity entity;

  UserModel({required String uid, required String email, String? displayName})
    : entity = UserEntity(uid: uid, email: email, displayName: displayName);

  String get uid => entity.uid;
  String get email => entity.email;
  String? get displayName => entity.displayName;

  factory UserModel.fromFirebaseUser(dynamic user) {
    return UserModel(
      uid: user.uid,
      email: user.email,
      displayName: user.displayName,
    );
  }

  UserEntity toEntity() => entity;
}
`,

	"auth_repository_impl.dart": `import 'package:firebase_auth/firebase_auth.dart';
import 'package:google_sign_in/google_sign_in.dart';
import '../../domain/entities/user_entity.dart';
import '../../domain/repositories/auth_repository.dart';
import '../models/user_model.dart';

class AuthRepositoryImpl implements AuthRepository {
  final FirebaseAuth _firebaseAuth;
  final GoogleSignIn _googleSignIn;

  AuthRepositoryImpl(this._firebaseAuth, this._googleSignIn);

  @override
  Future<UserEntity> signIn(String email, String password) async {
    final userCredential = await _firebaseAuth.signInWithEmailAndPassword(
      email: email,
      password: password,
    );

    final userModel = UserModel.fromFirebaseUser(userCredential.user);
    return userModel.toEntity();
  }

  @override
  Future<UserEntity> signUp(String name, String email, String password) async {
    final userCredential = await _firebaseAuth.createUserWithEmailAndPassword(
      email: email,
      password: password,
    );

    await userCredential.user?.updateDisplayName(name);

    final userModel = UserModel.fromFirebaseUser(userCredential.user);
    return userModel.toEntity();
  }

  @override
  Future<UserEntity> signInWithGoogle() async {
    try {
      final GoogleSignInAccount? googleUser = await _googleSignIn.signIn();

      if (googleUser == null) {
        throw Exception("Google sign in was canceled");
      }

      final GoogleSignInAuthentication googleAuth =
          await googleUser.authentication;
      final AuthCredential credential = GoogleAuthProvider.credential(
        accessToken: googleAuth.accessToken,
        idToken: googleAuth.idToken,
      );

      final UserCredential userCredential = await _firebaseAuth
          .signInWithCredential(credential);
      final userModel = UserModel.fromFirebaseUser(userCredential.user);
      return userModel.toEntity();
    } catch (e) {
      throw Exception("Failed to sign in with Google: ${e.toString()}");
    }
  }

  @override
  Future<void> signOut() async {
    await Future.wait([_firebaseAuth.signOut(), _googleSignIn.signOut()]);
  }

  @override
  Future<UserEntity?> getUser() async {
    try {
      final user = _firebaseAuth.currentUser;
      if (user != null) {
        final userModel = UserModel.fromFirebaseUser(user);
        return userModel.toEntity();
      }
      return null;
    } catch (e) {
      throw Exception("Failed to get user: ${e.toString()}");
    }
  }
}
`,

	"user_entity.dart": `import 'package:freezed_annotation/freezed_annotation.dart';

part 'user_entity.g.dart';
part 'user_entity.freezed.dart';

@freezed
abstract class UserEntity with _$UserEntity {
  const UserEntity._();
  const factory UserEntity({
    required String uid,
    required String email,
    String? displayName,
  }) = _UserEntity;

  factory UserEntity.fromJson(Map<String, dynamic> json) =>
      _$UserEntityFromJson(json);
}`,

	"auth_repository.dart": `import '../entities/user_entity.dart';

abstract class AuthRepository {
  Future<UserEntity> signIn(String email, String password);
  Future<UserEntity> signUp(String name, String email, String password);
  Future<UserEntity> signInWithGoogle();
  Future<void> signOut();
  Future<UserEntity?> getUser();
}`,

	"google_signin_usecase.dart": `import '../entities/user_entity.dart';
import '../repositories/auth_repository.dart';

class GoogleSignInUseCase {
  final AuthRepository repository;

  GoogleSignInUseCase(this.repository);

  Future<UserEntity> execute() async {
    return await repository.signInWithGoogle();
  }
}
`,

	"signin_usecase.dart": `import '../entities/user_entity.dart';
import '../repositories/auth_repository.dart';

class SignInUseCase {
  final AuthRepository repository;

  SignInUseCase(this.repository);

  Future<UserEntity> execute(String email, String password) async {
    return await repository.signIn(email, password);
  }
}
`,

	"signup_usecase.dart": `import '../entities/user_entity.dart';
import '../repositories/auth_repository.dart';

class SignUpUseCase {
  final AuthRepository repository;

  SignUpUseCase(this.repository);

  Future<UserEntity> execute(String name, String email, String password) async {
    return await repository.signUp(name, email, password);
  }
}
`,

	"signout_usecase.dart": `import '../repositories/auth_repository.dart';

class SignOutUseCase {
  final AuthRepository repository;

  SignOutUseCase(this.repository);

  Future<void> execute() async {
    await repository.signOut();
  }
}
`,

	"auth_bloc.dart": `import '/features/auth/domain/repositories/auth_repository.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'auth_event.dart';
import 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository authRepository;

  AuthBloc(this.authRepository) : super(const AuthState.initial()) {
    on<AuthEvent>((event, emit) async {
      final eventString = event.toString();

      if (eventString.contains('signIn(')) {
        final signInEvent = event as dynamic;
        emit(const AuthState.loading());
        try {
          final user = await authRepository.signIn(
            signInEvent.email,
            signInEvent.password,
          );
          emit(AuthState.authenticated(user));
        } catch (e) {
          emit(AuthState.error(e.toString()));
        }
      } else if (eventString.contains('signUp')) {
        final signUpEvent = event as dynamic;
        emit(const AuthState.loading());
        try {
          final user = await authRepository.signUp(
            signUpEvent.name,
            signUpEvent.email,
            signUpEvent.password,
          );
          emit(AuthState.authenticated(user));
        } catch (e) {
          emit(AuthState.error(e.toString()));
        }
      } else if (eventString.contains('signInWithGoogle')) {
        emit(const AuthState.loading());
        try {
          final user = await authRepository.signInWithGoogle();
          emit(AuthState.authenticated(user));
        } catch (e) {
          emit(AuthState.error(e.toString()));
        }
      } else if (eventString.contains('signOut')) {
        await authRepository.signOut();
        emit(const AuthState.unauthenticated());
      } else if (eventString.contains('checkAuthStatus')) {
        final user = await authRepository.getUser();
        user != null
            ? emit(AuthState.authenticated(user))
            : emit(const AuthState.unauthenticated());
      }
    });
  }
}
`,

	"auth_event.dart": `import 'package:freezed_annotation/freezed_annotation.dart';

part 'auth_event.freezed.dart';

@freezed
class AuthEvent with _$AuthEvent {
  const factory AuthEvent.signIn(String email, String password) = _SignIn;
  const factory AuthEvent.signUp(String name, String email, String password) =
      _SignUp;
  const factory AuthEvent.signInWithGoogle() = _SignInWithGoogle;
  const factory AuthEvent.signOut() = _SignOut;
  const factory AuthEvent.checkAuthStatus() = _CheckAuthStatus;
}
`,

	"auth_state.dart": `import '/features/auth/domain/entities/user_entity.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'auth_state.freezed.dart';

@freezed
class AuthState with _$AuthState {
  const factory AuthState.initial() = _Initial;
  const factory AuthState.loading() = _Loading;
  const factory AuthState.authenticated(UserEntity user) = _Authenticated;
  const factory AuthState.unauthenticated() = _Unauthenticated;
  const factory AuthState.error(String message) = _Error;
}
`,

	"signin_page.dart": `import '/core/constants/routes.dart';
import 'package:flutter/material.dart';
import 'package:shadcn_ui/shadcn_ui.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '/features/auth/presentation/bloc/auth_bloc.dart';
import '/features/auth/presentation/bloc/auth_event.dart';
import '/features/auth/presentation/bloc/auth_state.dart';

class SignInScreen extends StatefulWidget {
  const SignInScreen({super.key});

  @override
  State<SignInScreen> createState() => _SignUpScreenState();
}

class _SignUpScreenState extends State<SignInScreen> {
  final formKey = GlobalKey<ShadFormState>();
  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  bool isLoading = false;

  @override
  void dispose() {
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  void _handleSignIn() {
    if (formKey.currentState?.validate() ?? false) {
      context.read<AuthBloc>().add(
        AuthEvent.signIn(emailController.text.trim(), passwordController.text),
      );
    }
  }

  void _handleGoogleSignIn() {
    context.read<AuthBloc>().add(const AuthEvent.signInWithGoogle());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        width: MediaQuery.of(context).size.width,
        height: MediaQuery.of(context).size.height,
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [Colors.purple.shade50, Colors.blue.shade50],
            stops: const [0.0, 1.0],
          ),
        ),
        child: BlocConsumer<AuthBloc, AuthState>(
          listener: (context, state) {
            final stateString = state.toString();

            if (stateString.contains('initial')) {
              setState(() => isLoading = false);
            } else if (stateString.contains('loading')) {
              setState(() => isLoading = true);
            } else if (stateString.contains('authenticated')) {
              setState(() => isLoading = false);
              AppRoutes.navigateToAndRemoveUntil(context, AppRoutes.home);
            } else if (stateString.contains('unauthenticated')) {
              setState(() => isLoading = false);
            } else if (stateString.contains('error')) {
              setState(() => isLoading = false);

              final errorMessage =
                  stateString.split('message: ')[1].split(')')[0];
              ScaffoldMessenger.of(
                context,
              ).showSnackBar(SnackBar(content: Text(errorMessage)));
            }
          },
          builder: (context, state) {
            return Padding(
              padding: const EdgeInsets.only(top: 150, left: 5),
              child: Column(
                children: [
                  Text(
                    'Get Started',
                    style: ShadTheme.of(context).textTheme.h1Large,
                  ),
                  Text(
                    'with F3 Stack',
                    style: ShadTheme.of(context).textTheme.h1Large,
                  ),
                  const SizedBox(height: 80),
                  ShadForm(
                    key: formKey,
                    child: ConstrainedBox(
                      constraints: const BoxConstraints(
                        maxWidth: 350,
                        maxHeight: 500,
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          ShadInputFormField(
                            id: 'email',
                            controller: emailController,
                            label: const Text('Email'),
                            placeholder: const Text('Enter your email'),
                            validator: (v) {
                              final emailRegex = RegExp(
                                r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$',
                              );
                              if (!emailRegex.hasMatch(v)) {
                                return 'Please enter a valid email address.';
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 16),
                          ShadInputFormField(
                            id: 'password',
                            controller: passwordController,
                            label: const Text('Password'),
                            placeholder: const Text('Enter your password'),
                            obscureText: true,
                            validator: (v) {
                              if (v.length < 6) {
                                return 'Password must be at least 6 characters.';
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 15),
                          Center(
                            child: Column(
                              children: [
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Text(
                                      'Dont have an account ?',
                                      style:
                                          ShadTheme.of(context).textTheme.muted,
                                    ),
                                    const SizedBox(width: 2),
                                    GestureDetector(
                                      onTap: () {
                                        Navigator.pushNamed(
                                          context,
                                          '/sign-up',
                                        );
                                      },
                                      child: Text(
                                        'Sign Up',
                                        style: ShadTheme.of(
                                          context,
                                        ).textTheme.small.copyWith(
                                          color: Theme.of(context).primaryColor,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                                const SizedBox(height: 5),
                                ShadButton(
                                  width: 200,
                                  onPressed: isLoading ? null : _handleSignIn,
                                  child:
                                      isLoading
                                          ? const SizedBox(
                                            width: 20,
                                            height: 20,
                                            child: CircularProgressIndicator(
                                              strokeWidth: 2,
                                            ),
                                          )
                                          : const Text('SIGN IN'),
                                ),
                                const SizedBox(height: 5),
                                Row(
                                  children: [
                                    const Expanded(
                                      child: Divider(
                                        color: Colors.grey,
                                        thickness: 0.5,
                                      ),
                                    ),
                                    Padding(
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 10,
                                      ),
                                      child: Text(
                                        'or with',
                                        style:
                                            ShadTheme.of(
                                              context,
                                            ).textTheme.muted,
                                      ),
                                    ),
                                    const Expanded(
                                      child: Divider(
                                        color: Colors.grey,
                                        thickness: 0.5,
                                      ),
                                    ),
                                  ],
                                ),
                                const SizedBox(height: 20),
                                ShadButton.outline(
                                  leading: Image.asset(
                                    'assets/images/google-icon.webp',
                                    scale: 30,
                                  ),
                                  width: 200,
                                  onPressed:
                                      isLoading ? null : _handleGoogleSignIn,
                                  child:
                                      isLoading
                                          ? const SizedBox(
                                            width: 20,
                                            height: 20,
                                            child: CircularProgressIndicator(
                                              strokeWidth: 2,
                                            ),
                                          )
                                          : const Text('Sign in with Google'),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 180),
                  Text(
                    'This Signin page is powered by Firebase Auth.\nSign in using Email & Password / Google Auth.',
                    style: ShadTheme.of(context).textTheme.muted,
                  ),
                ],
              ),
            );
          },
        ),
      ),
    );
  }
}
`,

	"signup_page.dart": `import '/core/constants/routes.dart';
import 'package:flutter/material.dart';
import 'package:shadcn_ui/shadcn_ui.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '/features/auth/presentation/bloc/auth_bloc.dart';
import '/features/auth/presentation/bloc/auth_event.dart';
import '/features/auth/presentation/bloc/auth_state.dart';

class SignUpScreen extends StatefulWidget {
  const SignUpScreen({super.key});

  @override
  State<SignUpScreen> createState() => _SignUpScreenState();
}

class _SignUpScreenState extends State<SignUpScreen> {
  final formKey = GlobalKey<ShadFormState>();
  final TextEditingController nameController = TextEditingController();
  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  final TextEditingController confirmPasswordController =
      TextEditingController();
  bool isLoading = false;

  @override
  void dispose() {
    nameController.dispose();
    emailController.dispose();
    passwordController.dispose();
    confirmPasswordController.dispose();
    super.dispose();
  }

  void _handleSignUp() {
    if (formKey.currentState?.validate() ?? false) {
      context.read<AuthBloc>().add(
        AuthEvent.signUp(
          nameController.text.trim(),
          emailController.text.trim(),
          passwordController.text,
        ),
      );
    }
  }

  void _handleGoogleSignIn() {
    context.read<AuthBloc>().add(const AuthEvent.signInWithGoogle());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        width: MediaQuery.of(context).size.width,
        height: MediaQuery.of(context).size.height,
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [Colors.purple.shade50, Colors.blue.shade50],
            stops: const [0.0, 1.0],
          ),
        ),
        child: BlocConsumer<AuthBloc, AuthState>(
          listener: (context, state) {
            final stateString = state.toString();

            if (stateString.contains('initial')) {
              setState(() => isLoading = false);
            } else if (stateString.contains('loading')) {
              setState(() => isLoading = true);
            } else if (stateString.contains('authenticated')) {
              setState(() => isLoading = false);
              AppRoutes.navigateToAndRemoveUntil(context, AppRoutes.home);
            } else if (stateString.contains('unauthenticated')) {
              setState(() => isLoading = false);
            } else if (stateString.contains('error')) {
              setState(() => isLoading = false);

              final errorMessage =
                  stateString.split('message: ')[1].split(')')[0];
              ScaffoldMessenger.of(
                context,
              ).showSnackBar(SnackBar(content: Text(errorMessage)));
            }
          },
          builder: (context, state) {
            return Padding(
              padding: const EdgeInsets.only(top: 100, left: 5),
              child: Column(
                children: [
                  Text(
                    'Create Account',
                    style: ShadTheme.of(context).textTheme.h1Large,
                  ),
                  Text(
                    'with F3 Stack',
                    style: ShadTheme.of(context).textTheme.h1Large,
                  ),
                  const SizedBox(height: 80),
                  Padding(
                    padding: const EdgeInsets.only(right: 10),
                    child: ShadForm(
                      key: formKey,
                      child: ConstrainedBox(
                        constraints: const BoxConstraints(
                          maxWidth: 350,
                          maxHeight: 500,
                        ),
                        child: Column(
                          children: [
                            ShadInputFormField(
                              id: 'name',
                              controller: nameController,
                              label: const Text('Full Name'),
                              placeholder: const Text('Enter your full name'),
                              validator: (v) {
                                if (v.isEmpty) {
                                  return 'Please enter your name.';
                                }
                                return null;
                              },
                            ),
                            const SizedBox(height: 16),
                            ShadInputFormField(
                              id: 'email',
                              controller: emailController,
                              label: const Text('Email'),
                              placeholder: const Text('Enter your email'),
                              validator: (v) {
                                final emailRegex = RegExp(
                                  r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$',
                                );
                                if (!emailRegex.hasMatch(v)) {
                                  return 'Please enter a valid email address.';
                                }
                                return null;
                              },
                            ),
                            const SizedBox(height: 16),
                            ShadInputFormField(
                              id: 'password',
                              controller: passwordController,
                              label: const Text('Password'),
                              placeholder: const Text('Enter your password'),
                              obscureText: true,
                              validator: (v) {
                                if (v.length < 6) {
                                  return 'Password must be at least 6 characters.';
                                }
                                return null;
                              },
                            ),
                            const SizedBox(height: 16),
                            ShadInputFormField(
                              id: 'confirmPassword',
                              controller: confirmPasswordController,
                              label: const Text('Confirm Password'),
                              placeholder: const Text('Confirm your password'),
                              obscureText: true,
                              validator: (v) {
                                if (v != passwordController.text) {
                                  return 'Passwords do not match.';
                                }
                                return null;
                              },
                            ),
                            const SizedBox(height: 15),
                            Center(
                              child: Column(
                                children: [
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Text(
                                        'Already have an account?',
                                        style:
                                            ShadTheme.of(
                                              context,
                                            ).textTheme.muted,
                                      ),
                                      const SizedBox(width: 2),
                                      GestureDetector(
                                        onTap: () {
                                          Navigator.pushReplacementNamed(
                                            context,
                                            '/sign-in',
                                          );
                                        },
                                        child: Text(
                                          'Sign In',
                                          style: ShadTheme.of(
                                            context,
                                          ).textTheme.small.copyWith(
                                            color:
                                                Theme.of(context).primaryColor,
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                  const SizedBox(height: 3),
                                  ShadButton(
                                    width: 200,
                                    onPressed: isLoading ? null : _handleSignUp,
                                    child:
                                        isLoading
                                            ? const SizedBox(
                                              width: 20,
                                              height: 20,
                                              child: CircularProgressIndicator(
                                                strokeWidth: 2,
                                              ),
                                            )
                                            : const Text('SIGN UP'),
                                  ),
                                  Row(
                                    children: [
                                      const Expanded(
                                        child: Divider(
                                          color: Colors.grey,
                                          thickness: 0.5,
                                        ),
                                      ),
                                      Padding(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: 10,
                                        ),
                                        child: Text(
                                          'or with',
                                          style:
                                              ShadTheme.of(
                                                context,
                                              ).textTheme.muted,
                                        ),
                                      ),
                                      const Expanded(
                                        child: Divider(
                                          color: Colors.grey,
                                          thickness: 0.5,
                                        ),
                                      ),
                                    ],
                                  ),
                                  ShadButton.outline(
                                    leading: Image.asset(
                                      'assets/images/google-icon.webp',
                                      scale: 30,
                                    ),
                                    width: 200,
                                    onPressed:
                                        isLoading ? null : _handleGoogleSignIn,
                                    child:
                                        isLoading
                                            ? const SizedBox(
                                              width: 20,
                                              height: 20,
                                              child: CircularProgressIndicator(
                                                strokeWidth: 2,
                                              ),
                                            )
                                            : const Text('Sign in with Google'),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(height: 70),
                  Text(
                    'This Signup page is powered by Firebase Auth.\nSign up using Email & Password / Google Auth.',
                    style: ShadTheme.of(context).textTheme.muted,
                  ),
                ],
              ),
            );
          },
        ),
      ),
    );
  }
}
`,

	"auth_wrapper.dart": `import '/features/auth/presentation/bloc/auth_bloc.dart';
import '/features/auth/presentation/bloc/auth_event.dart';
import '/features/auth/presentation/bloc/auth_state.dart';
import 'signin_page.dart';
import '../../../home/presentation/pages/home_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class AuthWrapper extends StatelessWidget {
  const AuthWrapper({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AuthBloc, AuthState>(
      builder: (context, state) {
        final stateString = state.toString();

        if (stateString == 'AuthState.initial()' ||
            stateString == 'AuthState.loading()') {
          return const Center(child: CircularProgressIndicator());
        } else if (stateString.startsWith('AuthState.authenticated')) {
          return const HomeScreen();
        } else if (stateString == 'AuthState.unauthenticated()') {
          return const SignInScreen();
        } else if (stateString.startsWith('AuthState.error')) {
          final errorMessage = _extractErrorMessage(stateString);
          return _buildErrorScreen(context, errorMessage);
        } else {
          return const Center(child: CircularProgressIndicator());
        }
      },
    );
  }

  String _extractErrorMessage(String stateString) {
    final regex = RegExp(r'message: (.*?)\)');
    final match = regex.firstMatch(stateString);
    if (match != null && match.groupCount >= 1) {
      return match.group(1) ?? 'Unknown error';
    }
    try {
      return stateString.split('message: ')[1].split(')')[0];
    } catch (e) {
      return 'Unknown error';
    }
  }

  Widget _buildErrorScreen(BuildContext context, String errorMessage) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('Error: $errorMessage'),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () {
                context.read<AuthBloc>().add(const AuthEvent.checkAuthStatus());
              },
              child: const Text('Retry'),
            ),
          ],
        ),
      ),
    );
  }
}
`,

	"home_page.dart": `import '../widgets/content.dart';
import 'package:flutter/material.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(body: const Content());
  }
}
`,

	"action_button.dart": `import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

class ActionButtons extends StatelessWidget {
  const ActionButtons({super.key});

  Future<void> _launch() async {
    const String githubUrl = 'https://github.com/mixter3011';
    final Uri url = Uri.parse(githubUrl);
    if (await canLaunchUrl(url)) {
      await launchUrl(url, mode: LaunchMode.externalApplication);
    } else {
      throw Exception('Could not launch $url');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        ElevatedButton.icon(
          onPressed: () {},
          icon: const Icon(Icons.arrow_forward, color: Colors.black),
          label: const Text(
            'Get Started',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: Colors.black,
            ),
          ),
          style: ElevatedButton.styleFrom(
            backgroundColor: Colors.white,
            padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 28),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        ),
        const SizedBox(width: 16),
        OutlinedButton.icon(
          onPressed: _launch,
          icon: Image.asset('assets/images/github-icon.png', scale: 20),
          label: const Text(
            'GitHub',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: Colors.white,
            ),
          ),
          style: OutlinedButton.styleFrom(
            backgroundColor: Colors.white.withOpacity(0.1),
            side: BorderSide(color: Colors.white.withOpacity(0.2)),
            padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 28),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        ),
      ],
    );
  }
}
`,

	"bottom_bar.dart": `import 'package:flutter/material.dart';

class CustomBottomNavBar extends StatelessWidget {
  const CustomBottomNavBar({super.key});

  @override
  Widget build(BuildContext context) {
    return BottomNavigationBar(
      items: const [
        BottomNavigationBarItem(
          icon: Icon(Icons.home_outlined),
          activeIcon: Icon(Icons.home),
          label: 'Home',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.book_outlined),
          activeIcon: Icon(Icons.book),
          label: 'Docs',
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.code),
          activeIcon: Icon(Icons.code),
          label: 'GitHub',
        ),
      ],
      currentIndex: 0,
      selectedItemColor: const Color(0xFF6C63FF),
      unselectedItemColor: Color(0xFF666666),
      backgroundColor: Colors.white,
      elevation: 8,
      onTap: (index) {
        // Navigation logic would go here
      },
    );
  }
}
`,

	"content.dart": `import 'hero.dart';
import 'features.dart';
import 'started.dart';
import 'package:flutter/material.dart';

class Content extends StatelessWidget {
  const Content({super.key});

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Column(children: const [HeroSec(), Features(), Started()]),
    );
  }
}
`,

	"feature_grid.dart": `import 'feature_card.dart';
import 'package:flutter/material.dart';

class FeatureGrid extends StatelessWidget {
  const FeatureGrid({super.key});

  @override
  Widget build(BuildContext context) {
    final List<Map<String, dynamic>> features = [
      {
        'icon': Icons.terminal,
        'title': 'Type Safety',
        'description':
            'Full stack type safety with Freezed. No more runtime errors.',
        'gradient': const [Color(0xFFFF6B6B), Color(0xFFFF8E8E)],
      },
      {
        'icon': Icons.bolt,
        'title': 'Performance',
        'description':
            'Built on Flutter for native performance on all platforms.',
        'gradient': const [Color(0xFF4ECDC4), Color(0xFF45B7AF)],
      },
      {
        'icon': Icons.shield,
        'title': 'Auth Built-in',
        'description':
            'Secure authentication with Firebase Auth, ready out of the box.',
        'gradient': const [Color(0xFFFFD93D), Color(0xFFF6C90E)],
      },
      {
        'icon': Icons.storage,
        'title': 'Database',
        'description': 'Powered by Firebase with real-time capabilities.',
        'gradient': const [Color(0xFF6C63FF), Color(0xFF5A52D5)],
      },
    ];

    return LayoutBuilder(
      builder: (context, constraints) {
        final isWideScreen = constraints.maxWidth > 800;

        return Wrap(
          spacing: 24,
          runSpacing: 24,
          children:
              features.map((feature) {
                return SizedBox(
                  width:
                      isWideScreen
                          ? (constraints.maxWidth / 2) - 36
                          : constraints.maxWidth,
                  child: FeatureCard(
                    icon: feature['icon'],
                    title: feature['title'],
                    description: feature['description'],
                    gradient: feature['gradient'],
                  ),
                );
              }).toList(),
        );
      },
    );
  }
}
`,

	"feature_card.dart": `import 'package:flutter/material.dart';
import 'package:shadcn_ui/shadcn_ui.dart';

class FeatureCard extends StatefulWidget {
  final IconData icon;
  final String title;
  final String description;
  final List<Color> gradient;

  const FeatureCard({
    super.key,
    required this.icon,
    required this.title,
    required this.description,
    required this.gradient,
  });

  @override
  State<FeatureCard> createState() => _FeatureCardState();
}

class _FeatureCardState extends State<FeatureCard> {
  bool _isHovering = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovering = true),
      onExit: (_) => setState(() => _isHovering = false),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              color:
                  _isHovering
                      ? widget.gradient[0].withOpacity(0.5)
                      : Colors.black.withOpacity(0.2),
              blurRadius: _isHovering ? 45 : 10,
              spreadRadius: _isHovering ? 2 : 0,
              offset: const Offset(0, 8),
            ),
          ],
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(16),
          child: Container(
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: widget.gradient,
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
            ),
            padding: const EdgeInsets.all(1),
            child: Container(
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(15),
              ),
              padding: const EdgeInsets.all(24),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  AnimatedContainer(
                    duration: const Duration(milliseconds: 200),
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(12),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.1),
                          blurRadius: 4,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: Icon(
                      widget.icon,
                      color:
                          _isHovering ? widget.gradient[1] : widget.gradient[0],
                      size: 24,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(widget.title, style: ShadTheme.of(context).textTheme.h3),
                  const SizedBox(height: 8),
                  Text(
                    widget.description,
                    style: ShadTheme.of(context).textTheme.lead,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
`,

	"features.dart": `import 'feature_grid.dart';
import 'package:flutter/material.dart';
import 'package:shadcn_ui/shadcn_ui.dart';

class Features extends StatelessWidget {
  const Features({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(40),
      child: Column(
        children: [
          Text(
            'Why F3?',
            style: ShadTheme.of(context).textTheme.h1,
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Text(
              'F3 brings together the best tools in the Flutter Ecosystem.',
              style: ShadTheme.of(context).textTheme.lead,
              textAlign: TextAlign.center,
            ),
          ),
          const SizedBox(height: 48),
          const FeatureGrid(),
        ],
      ),
    );
  }
}
`,

	"hero.dart": `import '/features/auth/presentation/bloc/auth_bloc.dart';
import '/features/auth/presentation/bloc/auth_event.dart';
import '/features/home/presentation/widgets/action_button.dart';

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:shadcn_ui/shadcn_ui.dart';

class HeroSec extends StatelessWidget {
  const HeroSec({super.key});

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Container(
          padding: const EdgeInsets.all(20),
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [Color(0xFF000B18), Color(0xFF0F2027), Color(0xFF203A43)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
            borderRadius: BorderRadius.only(
              bottomLeft: Radius.circular(40),
              bottomRight: Radius.circular(40),
            ),
          ),
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.only(top: 30),
                child: Align(
                  alignment: Alignment.topLeft,
                  child: ShadButton(
                    leading: const Icon(LucideIcons.logOut),
                    size: ShadButtonSize.sm,
                    onPressed: () {
                      context.read<AuthBloc>().add(const AuthEvent.signOut());
                    },
                    child: Text('Log Out'),
                  ),
                ),
              ),
              Image.asset('assets/images/logo.png', scale: 5),
              const SizedBox(height: 16),
              Text(
                'F3 A Modern Flutter Framework',
                style: GoogleFonts.inter(
                  color: Colors.white,
                  fontSize: 36,
                  fontWeight: FontWeight.bold,
                  textStyle: ShadTheme.of(context).textTheme.h1Large,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Text(
                  'The modern typesafe Flutter stack for building cross-platform apps',
                  style: TextStyle(
                    fontSize: 18,
                    color: Colors.white.withOpacity(0.8),
                    height: 1.3,
                  ),
                  textAlign: TextAlign.center,
                ),
              ),
              const SizedBox(height: 32),
              const ActionButtons(),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ],
    );
  }
}
`,

	"started.dart": `import 'package:flutter/material.dart';
import 'package:shadcn_ui/shadcn_ui.dart';

class Started extends StatelessWidget {
  const Started({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(40),
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          colors: [Color(0xFFF8F9FA), Color(0xFFE9ECEF)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
      ),
      child: Column(
        children: [
          Text(
            'Get Started',
            style: ShadTheme.of(context).textTheme.h1,
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 10),
          Text(
            'Start by modifying this project  Comes fully with the F3 Stack.',
            style: ShadTheme.of(context).textTheme.lead,
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          Container(
            width: double.infinity,
            constraints: const BoxConstraints(maxWidth: 400),
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(12),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.1),
                  blurRadius: 4,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: const Text(
              'create-f3-app@latest',
              style: TextStyle(
                fontFamily: 'monospace',
                fontSize: 16,
                color: Colors.black,
              ),
            ),
          ),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: () {},
            icon: const Icon(Icons.open_in_new, color: Color(0xFF6C63FF)),
            label: const Text(
              'Read the docs',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Color(0xFF6C63FF),
              ),
            ),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.white,
              foregroundColor: const Color(0xFF6C63FF),
              padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 16),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              elevation: 2,
            ),
          ),
        ],
      ),
    );
  }
}
`,

	"Info.plist": `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>GIDClientID</key>
	<!-- Copied from GoogleService-Info.plist key CLIENT_ID -->
	<string>[YOUR IOS CLIENT ID]</string>
	<key>CFBundleURLTypes</key>
	<array>
		<dict>
			<key>CFBundleTypeRole</key>
			<string>Editor</string>
			<key>CFBundleURLSchemes</key>
			<array>
				<!-- Copied from GoogleService-Info.plist key REVERSED_CLIENT_ID -->
				<string>com.googleusercontent.apps.861823949799-vc35cprkp249096uujjn0vvnmcvjppkn</string>
			</array>
		</dict>
	</array>
	<key>CFBundleDevelopmentRegion</key>
	<string>$(DEVELOPMENT_LANGUAGE)</string>
	<key>CFBundleDisplayName</key>
	<string>Authf</string>
	<key>CFBundleExecutable</key>
	<string>$(EXECUTABLE_NAME)</string>
	<key>CFBundleIdentifier</key>
	<string>$(PRODUCT_BUNDLE_IDENTIFIER)</string>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleName</key>
	<string>authf</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>$(FLUTTER_BUILD_NAME)</string>
	<key>CFBundleSignature</key>
	<string>????</string>
	<key>CFBundleVersion</key>
	<string>$(FLUTTER_BUILD_NUMBER)</string>
	<key>LSRequiresIPhoneOS</key>
	<true/>
	<key>UILaunchStoryboardName</key>
	<string>LaunchScreen</string>
	<key>UIMainStoryboardFile</key>
	<string>Main</string>
	<key>UISupportedInterfaceOrientations</key>
	<array>
		<string>UIInterfaceOrientationPortrait</string>
		<string>UIInterfaceOrientationLandscapeLeft</string>
		<string>UIInterfaceOrientationLandscapeRight</string>
	</array>
	<key>UISupportedInterfaceOrientations~ipad</key>
	<array>
		<string>UIInterfaceOrientationPortrait</string>
		<string>UIInterfaceOrientationPortraitUpsideDown</string>
		<string>UIInterfaceOrientationLandscapeLeft</string>
		<string>UIInterfaceOrientationLandscapeRight</string>
	</array>
	<key>CADisableMinimumFrameDurationOnPhone</key>
	<true/>
	<key>UIApplicationSupportsIndirectInputEvents</key>
	<true/>
</dict>
</plist>
`,

	"pubspec.yaml": `name: f3stack
description: "A new Flutter project."

publish_to: 'none' 

version: 1.0.0+1

environment:
  sdk: ^3.7.0
  
dependencies:
  flutter:
    sdk: flutter

  cupertino_icons: ^1.0.8
  flutter_bloc: ^9.1.0
  freezed_annotation: ^3.0.0
  firebase_core: ^3.13.0
  firebase_auth: ^5.5.2
  json_annotation: ^4.9.0
  equatable: ^2.0.7
  shadcn_ui: ^0.24.0
  google_fonts: ^6.2.1
  url_launcher: ^6.3.1
  cached_network_image: ^3.4.1
  url_launcher_ios: ^6.3.3
  google_sign_in: ^6.3.0

dependency_overrides:
  analyzer: ^6.2.0

dev_dependencies:
  flutter_test:
    sdk: flutter

  flutter_lints: ^5.0.0
  build_runner: ^2.4.0
  freezed: ^3.0.6
  json_serializable: ^6.7.0

# The following section is specific to Flutter packages.
flutter:
  uses-material-design: true

  # To add assets to your application, add an assets section, like this:
  assets:
    - assets/images/
`,
}
