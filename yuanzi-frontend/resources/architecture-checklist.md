# Architecture Checklist

## Component Structure
- [ ] Component is modular and focused on single responsibility
- [ ] Props interface is defined with TypeScript
- [ ] Component name follows PascalCase convention
- [ ] File name matches component name

## Styling
- [ ] Uses Tailwind CSS classes (no inline styles)
- [ ] Uses theme colors (not hardcoded hex values)
- [ ] Responsive design considerations
- [ ] Consistent spacing (using theme spacing scale)

## Code Quality
- [ ] No console.log statements
- [ ] Proper error handling
- [ ] Accessible (ARIA attributes where needed)
- [ ] Comments for complex logic

## Performance
- [ ] React.memo used for expensive components
- [ ] useCallback/useMemo for handlers and computed values
- [ ] Lazy loading for heavy components

## Data Management
- [ ] Static data extracted to mockData.ts
- [ ] API calls use custom hooks
- [ ] Loading and error states handled

## Testing Readiness
- [ ] Component is testable (pure functions)
- [ ] Props are properly typed
- [ ] Edge cases considered
