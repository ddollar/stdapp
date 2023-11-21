declare const _default: __VLS_WithTemplateSlots<import("vue").DefineComponent<{
    action: {
        type: FunctionConstructor;
        required: true;
    };
    class: {
        type: StringConstructor;
        default: string;
    };
}, {}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
    show: (...args: any[]) => void;
    shown: (...args: any[]) => void;
}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    action: {
        type: FunctionConstructor;
        required: true;
    };
    class: {
        type: StringConstructor;
        default: string;
    };
}>> & {
    onShow?: ((...args: any[]) => any) | undefined;
    onShown?: ((...args: any[]) => any) | undefined;
}, {
    class: string;
}, {}>, {
    link?(_: {}): any;
    title?(_: {}): any;
    body?(_: {}): any;
    action?(_: {}): any;
}>;
export default _default;
type __VLS_WithTemplateSlots<T, S> = T & {
    new (): {
        $slots: S;
    };
};
