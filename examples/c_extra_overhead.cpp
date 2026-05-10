#include <bits/stdc++.h>
using namespace std;

int main() {
    int n;
    cin >> n;
    
    int cnt[5] = {0};
    for (int i = 0; i < n; i++) {
        int s; cin >> s;
        cnt[s]++;
    }
    
    int taxis = 0;
    
    taxis += cnt[4];
    
    taxis += cnt[3];
    int ones_used = min(cnt[1], cnt[3]);
    cnt[1] -= ones_used;
    
    taxis += cnt[2] / 2;
    if (cnt[2] % 2 == 1) {
        taxis++;
        int fill = min(cnt[1], 2);
        cnt[1] -= fill;
    }
    
    taxis += (cnt[1] + 3) / 4;
    
    cout << taxis << endl;
    return 0;
}
