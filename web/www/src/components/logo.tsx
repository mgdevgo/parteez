import { motion } from "framer-motion"

const animations = {
    start: {
        paddingTop: "50%",
        paddingBottom: "50%",
        width: "100%",
        opacity: 0,
    },
    mid: {
        opacity: 1,
        transition: { duration: 1, delay: 0.2 }
    },
    end: {
        paddingTop: 0,
        paddingBottom: 0,
        width: 164,
        transition: {
            duration: 1.1,
            delay: 1,
        }
    }
}

export default function Logo() {
    return (
        < motion.img
            src="/images/icon.png"
            alt="iditusi logo"
            initial="start"
            animate={["start", "mid", "end"]}
            variants={animations}
        />
    )
}