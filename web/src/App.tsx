import Logo from "@/components/logo"
import Copyright from "@/components/copyright"
import DownloadButton from "@/components/download-button"

export default function App() {
  return (
    <main className="min-h-dvh transition flex flex-col items-center">
      <div className="container flex justify-center pt-20 px-10">
        <div className="grid gap-6 items-start justify-items-center">

          <Logo />

          <div className="grid gap-3 justify-items-center">
            <h1 className="text-7xl font-druk font-extrabold hover:text-stroke-2 hover:text-stroke-black hover:text-stroke-fill-transparent">
              IDITUSI
            </h1>
            <h2 className="text-xl text-center">
              Лучшее приложение с тусовками
            </h2>
          </div>

          <DownloadButton />

        </div>
      </div>

      {/* Align to bottom */}
      <div className="flex mt-auto w-full"></div>

      <Copyright />
    </main>
  )
}