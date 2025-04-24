class Utils {
  String namesplite({required String name}) {
    if (name.length <= 14)
      return name;
    else {
      List<String> splited = name.split(" ");
      String res = "";
      for (int i = 0; i < splited.length; i++) {
        if (res.length + splited[i].length <= 14) {
          res += splited[i] + " ";
        } else
          return res;
      }
    }
    return name;
  }
}
