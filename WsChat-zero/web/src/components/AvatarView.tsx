import type { CSSProperties } from 'react';

type AvatarViewProps = {
  src?: string;
  alt: string;
  size?: number;
  radius?: number;
  className?: string;
  style?: CSSProperties;
  onClick?: () => void;
  title?: string;
};

export default function AvatarView({
  src,
  alt,
  size = 36,
  radius = 12,
  className,
  style,
  onClick,
  title,
}: AvatarViewProps) {
  const commonStyle: CSSProperties = {
    width: size,
    height: size,
    borderRadius: radius,
    flexShrink: 0,
    objectFit: 'cover',
    background: '#eef3fb',
    cursor: onClick ? 'pointer' : 'default',
    ...style,
  };

  if (src) {
    return <img src={src} alt={alt} title={title} className={className} style={commonStyle} onClick={onClick} />;
  }

  return <div aria-label={alt} title={title} className={className} style={commonStyle} onClick={onClick} />;
}
