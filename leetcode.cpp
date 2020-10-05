class Solution {
    public List<String> buildArray(int[] target, int n) {
        List<String> res = new ArrayList<>();
        int current = 0;
        int p = target.length;
        for (int i = 1; i <= n; i++) {
            if (current < p) {
                if (target[current] == i) {
                    res.add("Push");
                    current++;
                } else {
                    res.add("Push");
                    res.add("Pop");
                }
            } else {
                return res;
            }
        }
        return res;
    }
}
