const fs = require('fs');
const input = fs.readFileSync(0, 'utf8').trim();
const num = parseInt(input, 10);

if (!isNaN(num)) {
    if (num % 2 === 0) {
        console.log("Yes\nYes");
    } else {
        console.log("No")
    }
}