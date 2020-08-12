package com.example.running_plus;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.NotificationCompat;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Build;
import android.os.Bundle;
import android.view.WindowManager;
import android.webkit.WebChromeClient;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.widget.Toast;

import com.google.firebase.iid.FirebaseInstanceId;
import com.google.firebase.messaging.FirebaseMessaging;

public class MainActivity extends AppCompatActivity {

    private long backPressedTime = 0;
    private final long FINISH_INTERVAL_TIME = 2000;

    private WebView mWebView;
    private final String urladdr = "https://www.runningplus.net/Main/Index.nm";
    private final String title = "타이틀";
    private final String content = "내용";

    Intent foregroundServiceIntent;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        getWindow().setFlags(WindowManager.LayoutParams.FLAG_FULLSCREEN,
                WindowManager.LayoutParams.FLAG_FULLSCREEN);
        setContentView(R.layout.activity_main);

        /*
        memory over run..
        if want only once subscribed remark(//) the subscribeToTopic(~) line
        and go onNewToken of MyFirebaseMessagingService.class
         */
        FirebaseMessaging.getInstance().subscribeToTopic("running_plus_fcm");
        FirebaseInstanceId.getInstance().getToken();


        // xml web page view
        mWebView = findViewById(R.id.testWebview);

        mWebView.getSettings().setUseWideViewPort(true);
        mWebView.getSettings().setJavaScriptEnabled(true);
        mWebView.getSettings().setLoadWithOverviewMode(true);

        mWebView.setWebChromeClient(new WebChromeClient());
        mWebView.setWebViewClient(new WebViewClient());
        mWebView.loadUrl(urladdr);


        if (null == UndeadService.serviceIntent) {
            foregroundServiceIntent = new Intent(this, UndeadService.class);
            startService(foregroundServiceIntent);

        } else {
            foregroundServiceIntent = UndeadService.serviceIntent;
        }
    }

    // quit process check
    public void onBackPressed(){
        long tempTime = System.currentTimeMillis();
        long intervalTime = tempTime - backPressedTime;

        if(0 <= intervalTime && FINISH_INTERVAL_TIME >= intervalTime) {
            finishAffinity();
            System.runFinalization();
            System.exit(0);
        } else {
            backPressedTime = tempTime;
            Toast.makeText(getApplicationContext(), "종료를 원하시면 한번 더 누르세요.", Toast.LENGTH_SHORT).show();
        }
    }

    @Override
    protected  void onDestroy() {
        super.onDestroy();

        if (null != foregroundServiceIntent) {
            stopService(foregroundServiceIntent);
            foregroundServiceIntent = null;
        }
    }
}