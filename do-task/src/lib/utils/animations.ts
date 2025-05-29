import { cubicOut } from 'svelte/easing';
import type { TransitionConfig } from 'svelte/transition';

/**
 * Custom transition that combines fly and scale effects
 * Inspired by shadcn-svelte patterns
 */
export type FlyAndScaleParams = {
	y?: number;
	x?: number;
	start?: number;
	duration?: number;
};

export const flyAndScale = (
	node: Element,
	params: FlyAndScaleParams = { y: -8, x: 0, start: 0.95, duration: 150 }
): TransitionConfig => {
	const style = getComputedStyle(node);
	const transform = style.transform === 'none' ? '' : style.transform;

	const scaleConversion = (valueA: number, scaleA: [number, number], scaleB: [number, number]) => {
		const [minA, maxA] = scaleA;
		const [minB, maxB] = scaleB;

		const percentage = (valueA - minA) / (maxA - minA);
		const valueB = percentage * (maxB - minB) + minB;

		return valueB;
	};

	const styleToString = (style: Record<string, number | string | undefined>): string => {
		return Object.keys(style).reduce((str, key) => {
			if (style[key] === undefined) return str;
			return str + key + ':' + style[key] + ';';
		}, '');
	};

	return {
		duration: params.duration ?? 200,
		delay: 0,
		css: (t) => {
			const y = scaleConversion(t, [0, 1], [params.y ?? 5, 0]);
			const x = scaleConversion(t, [0, 1], [params.x ?? 0, 0]);
			const scale = scaleConversion(t, [0, 1], [params.start ?? 0.95, 1]);

			return styleToString({
				transform: transform + 'translate3d(' + x + 'px, ' + y + 'px, 0) scale(' + scale + ')',
				opacity: t
			});
		},
		easing: cubicOut
	};
};

/**
 * Staggered reveal animation for list items
 */
export const staggerReveal = (node: Element, { delay = 0 }: { delay?: number }) => {
	return {
		delay,
		duration: 400,
		css: (t: number) => {
			const eased = cubicOut(t);
			return `
				transform: translateY(${(1 - eased) * 20}px);
				opacity: ${eased};
			`;
		}
	};
};

/**
 * Smooth hover scale effect
 */
export const hoverScale = (node: HTMLElement, scale = 1.02) => {
	const handleMouseEnter = () => {
		node.style.transform = `scale(${scale})`;
		node.style.transition = 'transform 200ms ease-out';
	};

	const handleMouseLeave = () => {
		node.style.transform = 'scale(1)';
	};

	node.addEventListener('mouseenter', handleMouseEnter);
	node.addEventListener('mouseleave', handleMouseLeave);

	return {
		destroy() {
			node.removeEventListener('mouseenter', handleMouseEnter);
			node.removeEventListener('mouseleave', handleMouseLeave);
		}
	};
};

/**
 * Ripple effect for interactive elements
 */
export const ripple = (node: HTMLElement, { color = 'rgba(255, 255, 255, 0.6)' } = {}) => {
	const handleClick = (event: Event) => {
		const mouseEvent = event as MouseEvent;
		const rect = node.getBoundingClientRect();
		const size = Math.max(rect.width, rect.height);
		const x = mouseEvent.clientX - rect.left - size / 2;
		const y = mouseEvent.clientY - rect.top - size / 2;

		const ripple = document.createElement('span');
		ripple.style.cssText = `
			position: absolute;
			width: ${size}px;
			height: ${size}px;
			left: ${x}px;
			top: ${y}px;
			background: ${color};
			border-radius: 50%;
			transform: scale(0);
			animation: ripple-animation 600ms ease-out;
			pointer-events: none;
		`;

		node.appendChild(ripple);

		// Remove ripple after animation
		setTimeout(() => {
			ripple.remove();
		}, 600);
	};

	// Add ripple animation CSS if not already present
	if (!document.querySelector('#ripple-styles')) {
		const style = document.createElement('style');
		style.id = 'ripple-styles';
		style.textContent = `
			@keyframes ripple-animation {
				to {
					transform: scale(4);
					opacity: 0;
				}
			}
		`;
		document.head.appendChild(style);
	}

	node.addEventListener('click', handleClick);
	node.style.position = 'relative';
	node.style.overflow = 'hidden';

	return {
		destroy() {
			node.removeEventListener('click', handleClick);
		}
	};
};
