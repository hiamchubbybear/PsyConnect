import 'dart:convert';
import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';

class LocalDatabaseService {
  static final LocalDatabaseService _instance = LocalDatabaseService._internal();
  factory LocalDatabaseService() => _instance;
  LocalDatabaseService._internal();

  Database? _database;

  Future<Database> get database async {
    if (_database != null) return _database!;
    _database = await _initDatabase();
    return _database!;
  }

  Future<Database> _initDatabase() async {
    final dbPath = await getDatabasesPath();
    final pathString = join(dbPath, 'psyconnect_offline.db');

    return await openDatabase(
      pathString,
      version: 1,
      onCreate: (db, version) async {
        // 1. Table for caching newsfeed posts
        await db.execute('''
          CREATE TABLE cached_posts(
            id TEXT PRIMARY KEY,
            data TEXT,
            updatedAt INTEGER
          )
        ''');

        // 2. Table for caching chat messages/conversations
        await db.execute('''
          CREATE TABLE cached_chats(
            id TEXT PRIMARY KEY,
            data TEXT,
            updatedAt INTEGER
          )
        ''');

        // 3. Table for caching schedule sessions
        await db.execute('''
          CREATE TABLE cached_sessions(
            id TEXT PRIMARY KEY,
            data TEXT,
            updatedAt INTEGER
          )
        ''');
      },
    );
  }

  // --- Generic Helpers ---

  Future<void> cacheData(String table, String id, Map<String, dynamic> jsonMap) async {
    final db = await database;
    await db.insert(
      table,
      {
        'id': id,
        'data': jsonEncode(jsonMap),
        'updatedAt': DateTime.now().millisecondsSinceEpoch,
      },
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  Future<List<Map<String, dynamic>>> getCachedData(String table) async {
    final db = await database;
    final List<Map<String, dynamic>> maps = await db.query(table, orderBy: 'updatedAt DESC');

    return maps.map((m) {
      final String rawData = m['data'] as String;
      return jsonDecode(rawData) as Map<String, dynamic>;
    }).toList();
  }

  Future<void> clearCache(String table) async {
    final db = await database;
    await db.delete(table);
  }
}
