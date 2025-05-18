FROM openjdk:8

# 安装 Android SDK
RUN wget
RUN unzip commandlinetools-linux-6858069_latest.zip -d /sdk
RUN rm commandlinetools-linux-6858069_latest.zip

# 设置环境变量
ENV ANDROID_HOME=/sdk
ENV PATH=$PATH:$ANDROID_HOME/tools/bin:$ANDROID_HOME/platform-tools:$ANDROID_HOME/build-tools/28.0.3

# 下载所需的 Android SDK 组件
RUN yes | sdkmanager --licenses
RUN sdkmanager "platforms;android-28" "build-tools;28.0.3"

# 将你的项目复制到容器中
COPY . /app
WORKDIR /app

# 编译并打包 APK
RUN ./gradlew assembleRelease