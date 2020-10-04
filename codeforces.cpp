#include<bits/stdc++.h>
using namespace std;
 
int main(){
     long long int i,t,n[100],a,m,c,d,f;
     cin>>t>>m;
     for(i=0;i<t;i++)
     {
          cin>>n[i];
     }
     sort(n,n+t);
     c=0;
     d=t-1;
     while(m>0)
     {
          m-=n[d];
          d--;
          c++;
     }
     cout<<c;
     return 0;
}
