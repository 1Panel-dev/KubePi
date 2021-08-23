export function randomNum(n) {
  var rnd = "";
  for (var i = 0; i < n; i++) rnd += Math.floor(Math.random() * 10);
  return rnd;
}
