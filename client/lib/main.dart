import 'package:flutter/material.dart';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'dart:async';
void main() {
  runApp(const MyApp());
}



class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        // This is the theme of your application.
        //
        // TRY THIS: Try running your application with "flutter run". You'll see
        // the application has a purple toolbar. Then, without quitting the app,
        // try changing the seedColor in the colorScheme below to Colors.green
        // and then invoke "hot reload" (save your changes or press the "hot
        // reload" button in a Flutter-supported IDE, or press "r" if you used
        // the command line to start the app).
        //
        // Notice that the counter didn't reset back to zero; the application
        // state is not lost during the reload. To reset the state, use hot
        // restart instead.
        //
        // This works for code too, not just values: Most code changes can be
        // tested with just a hot reload.
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
      ),
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  Timer? _timer;
  String _response = 'Waiting for response...';
  int _counter = 0;
  
  
  Future<void> changeWallpaper(String imagePath) async {
  try {
    await Process.run('gsettings', [
      'set',
      'org.gnome.desktop.background',
      'picture-uri',
      'file://$imagePath'
    ]);
  } catch (e) {
    print('Error changing wallpaper: $e');
  }
}

Future<void> downloadAndSaveImage() async {
 const url = 'http://localhost:8080/api/wallpaper';
 const savePath = '/home/dev/background-sync-backgrounds/current.jpeg';
 
  try {
      final response = await http.get(Uri.parse(url));
      if (response.statusCode == 200) {
        final file = File(savePath);
        await file.writeAsBytes(response.bodyBytes);
        setState((){
          _response = "Wallpaper updated!";
        });
        print('saved');
      } else {
        print('Failed  ${response.statusCode}');
      }
    } catch (e) {
      print('Error: $e');
    }
 
 

}

Future<void> syncWallPaper() async{
 
    await downloadAndSaveImage();
    const imagePath ='/home/dev/background-sync-backgrounds/current.jpeg';
    changeWallpaper(imagePath);
}

void _startRequestLoop() {
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) async {
      try {
         await syncWallPaper();
       } catch (e) {
        setState(() {
          _response = 'Error: $e';
        });
      }
    });
  }


  @override
  void initState() {
    super.initState();
    _startRequestLoop();
  }

  void _onButtonPressed() {
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
          Text(
          _response,
          textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
     );
  }
}
