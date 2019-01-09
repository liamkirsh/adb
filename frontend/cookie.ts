// from http://stackoverflow.com/questions/10730362/get-cookie-by-name
export function getCookie(name: string) {
  var value = '; ' + document.cookie;
  var parts = value.split('; ' + name + '=');
  if (parts.length == 2) return parts[1].split(';').shift();
}

// from http://stackoverflow.com/questions/2144386/javascript-delete-cookie
export function deleteCookie(name: string) {
  document.cookie = name + '=;path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}
