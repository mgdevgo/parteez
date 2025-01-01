import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCopyright } from '@fortawesome/free-regular-svg-icons'

export default function Copyright() {
    var year = new Date().getFullYear()
    return (
        <div className="border-t border-black/5 w-full">
            <footer className="px-14 pt-8 pb-7">
                <div className="container">
                    <div className="text-xs font-mono font-medium">
                        <div className=" inline-flex items-center justify-center gap-1">
                            <FontAwesomeIcon icon={faCopyright} className='text-[10px]' />
                            <span>{year}</span>
                            <span>Max Gerasimov</span>
                        </div>
                        <p>www.iditusi.app v.0.1.1</p>
                    </div>
                </div>
            </footer>
        </div>
    )
}
