import "./globals.css";
import "./normalize.css";
import styles from "./Layout.module.css";
import Layouts from "./Layouts";

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
          <Layouts>{children}</Layouts>
        </div>
      </body>
    </html>
  );
}
