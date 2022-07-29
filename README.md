# Implementing Feature Flags With Go (RESt API)
Feature flags (feature toggles) adalah teknik yang digunakan untuk memodifikasi perilaku sistem tanpa harus merubah kode yang ada.
<blockquote cite="https://martinfowler.com/articles/feature-toggles.html">
  <p>
    <i>
      <h6>“Feature Toggles (often also refered to as Feature Flags) are a powerful technique,
           allowing teams to modify system behavior without changing code.”
       </h6>
    </i>
    — Martin Fowler</h5>
  <p>
</blockquote>
Seperti halnya saklar lampu, developer bisa mengaktifkan dan menonaktifkan suatu feature sesuai dengan kebutuhan. Teknik ini sangat bermanfaat ketika developer ingin release suatu feature baru untuk user atau proses yang memenuhi kriteria, hal ini digunakan untuk memastikan feature yang di release tidak ada kendala dan ketika ada kendala dapat langsung dimatikan tanpa harus merubah kode dan beralih semua proses menggunakan feature sebelumnya. Dan selain itu masih banyak hal yang dapat dilakukan dengan teknik feature flags.

Pada penerapannya, saya membuat sistem backend (Golang) yang menerapkan feature flags dimana package yang digunakan `go-feature-flag`. Sangat banyak kemudahan yang diberikan oleh package tersebut yang dapat dibaca pada [dokumentasi](https://docs.gofeatureflag.org/v0.27.1).

## Requirements
- Go version 1.18.x

## Dependencies
- [Echo web framework](https://echo.labstack.com)
- [Go dot env](https://github.com/joho/godotenv)
- [Go feature flag](https://github.com/thomaspoignant/go-feature-flag)

## Services
- RESt API
  - GET `/payment/:payment_id` </br>
    Mengambil pembayaran berdasarkan id, dengan penerapan `feature flagging` untuk role tertentu yang menggunakan feature baru dan selain itu akan 
    menggunakan existing feature.
  - GET `/generate` </br>
    Penerapan penyebaran feature baru dengan percentage yang di tentukan, percentage sendiri adalah suatu peluang user mendapatkan feature tsb.
    Untuk lebih jelasnya dapat melihat pada [dokumentasi](https://docs.gofeatureflag.org/v0.27.1/flag_format).
- Mock RESt API Notifier
  - POST `/notifier` </br>
    Digunakan untuk menerima pemberitahuan perubahan file config pada services yang menggunakan `feature flagging`, pada contoh ini menerapkan api-key
    pada header dengan model auth HMAC(SHA256)
## References
- [Pentingnya Feature Toggles Untuk Mobile Apps](https://medium.com/easyread/pentingnya-feature-toggles-feature-flags-untuk-mobile-apps-a31302c247f9)
- [Feature Toggles (aka Feature Flags)](https://martinfowler.com/articles/feature-toggles.html)
