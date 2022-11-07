import "./globals.css";
import "./normalize.css";
import styles from "./Layout.module.css";
import Template from "./template";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html>
      <head></head>
      <body>
        <div className={styles.app}>
          <Template>{children}</Template>
        </div>
      </body>
    </html>
  );
}
