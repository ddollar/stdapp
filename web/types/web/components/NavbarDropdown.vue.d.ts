declare const _default: __VLS_WithTemplateSlots<import("vue").DefineComponent<{
    icon: {
        type: StringConstructor;
        default: null;
    };
}, {}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
    click: (...args: any[]) => void;
}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    icon: {
        type: StringConstructor;
        default: null;
    };
}>> & {
    onClick?: ((...args: any[]) => any) | undefined;
}, {
    icon: string;
}, {}>, {
    default?(_: {}): any;
    items?(_: {}): any;
}>;
export default _default;
type __VLS_WithTemplateSlots<T, S> = T & {
    new (): {
        $slots: S;
    };
};
