import { useEffect, useState } from 'react';

const MOBILE_MAX_WIDTH = 768;

export function useMobile() {
  const [isMobile, setIsMobile] = useState(() => {
    if (typeof window === 'undefined') return false;
    return window.innerWidth <= MOBILE_MAX_WIDTH;
  });

  useEffect(() => {
    const onResize = () => setIsMobile(window.innerWidth <= MOBILE_MAX_WIDTH);
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, []);

  return isMobile;
}

export default useMobile;
